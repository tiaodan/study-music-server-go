package service

import (
	"fmt"
	"log"
	"strings"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
	"study-music-server-go/utils"
)

type RankService struct {
	rankMapper       *mapper.RankMapper
	singerRankMapper *mapper.SingerRankMapper
	albumRankMapper  *mapper.AlbumRankMapper
	songRankMapper   *mapper.SongRankMapper
	websiteMapper    *mapper.WebsiteMapper
}

func NewRankService() *RankService {
	return &RankService{
		rankMapper:       mapper.NewRankMapper(),
		singerRankMapper: mapper.NewSingerRankMapper(),
		albumRankMapper:  mapper.NewAlbumRankMapper(),
		songRankMapper:   mapper.NewSongRankMapper(),
		websiteMapper:    mapper.NewWebsiteMapper(),
	}
}

// getWebsiteName 获取网站名称
func (s *RankService) getWebsiteName(websiteId uint) string {
	websiteMap := map[uint]string{
		1: "qqmusic",
		2: "kugou",
		3: "kuwo",
		4: "netease",
		5: "migu",
	}
	if name, ok := websiteMap[websiteId]; ok {
		return name
	}
	return "unknown"
}

// buildFileNameIndex 构建文件名索引：去掉空格后key=歌手-歌名 -> 文件名
func buildFileNameIndex(files []utils.MusicFileInfo) map[string]string {
	index := make(map[string]string)
	for _, f := range files {
		key := strings.ReplaceAll(f.Singer, " ", "") + "-" + strings.ReplaceAll(f.SongName, " ", "")
		index[key] = f.OriginalName
	}
	return index
}

// ImportRank 导入排行榜数据
// 新逻辑：先入库歌曲（不管文件），再读取文件夹匹配更新 nas_url_path
func (s *RankService) ImportRank(req *models.RankImportRequest) *common.Response {
	log.Printf("开始导入排行榜: websiteId=%d, rankName=%s, folderPath=%s, 歌曲数量=%d",
		req.WebsiteId, req.RankName, req.FolderPath, len(req.List))

	websiteName := s.getWebsiteName(req.WebsiteId)

	var importedSongs []map[string]interface{}
	var failed []string
	var matchedFiles []string

	// 第一阶段：入库所有歌曲（不管文件是否存在）
	for _, item := range req.List {
		// 1. 处理歌手
		singerNames := strings.Split(item.Singer, "、")
		var singerIds []uint
		var firstSingerId uint
		var fullNameSinger string

		for i, singerName := range singerNames {
			singerName = strings.TrimSpace(singerName)
			if singerName == "" {
				continue
			}

			// 查询或创建歌手
			singers, _ := s.singerRankMapper.FindByName(singerName)
			var singerId uint
			if len(singers) == 0 {
				singer := &models.SingerRank{Name: singerName}
				if err := s.singerRankMapper.Add(singer); err != nil {
					failed = append(failed, fmt.Sprintf("创建歌手失败: %s - %v", singerName, err))
					continue
				}
				singerId = singer.ID
			} else {
				singerId = singers[0].ID
			}

			singerIds = append(singerIds, singerId)
			if i == 0 {
				firstSingerId = singerId
			}
		}

		if len(singerIds) == 0 {
			failed = append(failed, fmt.Sprintf("无法获取歌手ID: %s - %s", item.Name, item.Singer))
			continue
		}

		if len(singerIds) > 1 {
			fullNameSinger = strings.Join(singerNames, "、")
		}

		// 2. 处理专辑
		albumName := strings.TrimSpace(item.AlbumName)
		var albumId *uint
		if albumName != "" {
			foundAlbum, err := s.albumRankMapper.FindBySingerIdAndName(firstSingerId, albumName)
			if err == nil && foundAlbum != nil {
				albumId = &foundAlbum.ID
			} else {
				album := &models.AlbumRank{
					Name:     albumName,
					SingerId: firstSingerId,
				}
				if err := s.albumRankMapper.Add(album); err != nil {
					log.Printf("创建专辑失败: %s, err: %v", albumName, err)
				} else {
					albumId = &album.ID
				}
			}
		}

		// 3. 入库歌曲（nas_url_path 先为空）
		song := &models.SongRank{
			SingerId:       &firstSingerId,
			AlbumId:        albumId,
			Name:           strings.TrimSpace(item.Name),
			FullNameSinger: fullNameSinger,
			NasUrlPath:     "", // 先为空，后面匹配文件后更新
			SpiderUrl:      "", // 接口未返回，留空
		}

		if err := s.songRankMapper.Add(song); err != nil {
			failed = append(failed, fmt.Sprintf("创建歌曲失败: %s - %v", item.Name, err))
			continue
		}

		// 4. 写入排行榜记录
		rank := &models.Rank{
			WebsiteId: req.WebsiteId,
			Name:      req.RankName,
			SongId:    song.ID,
			AlbumId:   albumId,
			Album:     albumName,
			Singer:    item.Singer,
		}
		if err := s.rankMapper.Add(rank); err != nil {
			log.Printf("创建排行榜记录失败: %s - %v", item.Name, err)
		}

		importedSongs = append(importedSongs, map[string]interface{}{
			"song_id":  song.ID,
			"rank_id":  rank.ID,
			"name":     song.Name,
			"singer":   item.Singer,
			"album":    albumName,
			"file_key": strings.ReplaceAll(item.Singer, " ", "") + "-" + strings.ReplaceAll(item.Name, " ", ""),
		})
	}

	log.Printf("第一阶段入库完成: 入库 %d 首, 失败 %d 首", len(importedSongs), len(failed))

	// 第二阶段：读取文件夹，匹配并更新 nas_url_path
	folderPath := strings.TrimSpace(req.FolderPath)
	var matchedCount int

	if folderPath != "" {
		files, err := utils.GetMusicFiles(folderPath)
		if err != nil {
			log.Printf("读取目录失败: %v", err)
		} else {
			log.Printf("文件夹 %s 下有 %d 个音乐文件", folderPath, len(files))

			fileIndex := buildFileNameIndex(files)
			log.Printf("文件名索引构建完成，共 %d 个", len(fileIndex))

			// 遍历已入库的歌曲，匹配文件
			for i, songInfo := range importedSongs {
				fileKey := songInfo["file_key"].(string)
				fileName, found := fileIndex[fileKey]

				// 尝试只用第一个歌手匹配
				if !found && strings.Contains(songInfo["singer"].(string), "、") {
					firstSinger := strings.Split(songInfo["singer"].(string), "、")[0]
					altKey := strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "") + "-" + strings.ReplaceAll(songInfo["name"].(string), " ", "")
					fileName, found = fileIndex[altKey]
				}

				if found {
					songId := songInfo["song_id"].(uint)
					nasUrlPath := fmt.Sprintf("/rank/%s/%s/%s", websiteName, req.RankName, fileName)

					// 更新 nas_url_path
					if err := s.songRankMapper.UpdateNasUrlPath(songId, nasUrlPath); err != nil {
						log.Printf("更新 nas_url_path 失败: song_id=%d, err=%v", songId, err)
					} else {
						matchedCount++
						matchedFiles = append(matchedFiles, fileName)
						importedSongs[i]["nas_url_path"] = nasUrlPath
						log.Printf("匹配成功: %s -> %s", songInfo["name"], fileName)
					}
				}
			}
		}
	}

	log.Printf("第二阶段匹配完成: 匹配 %d 个文件", matchedCount)

	result := map[string]interface{}{
		"total":         len(req.List),
		"imported":      len(importedSongs),
		"failed":        len(failed),
		"matched_files": matchedCount,
		"results":       importedSongs,
		"failed_list":   failed,
		"matched_list":  matchedFiles,
	}

	log.Printf("排行榜导入完成: total=%d, imported=%d, failed=%d, matched=%d",
		len(req.List), len(importedSongs), len(failed), matchedCount)

	return common.SuccessWithData("导入完成", result)
}

// GetRankList 获取榜单列表
func (s *RankService) GetRankList(websiteId uint) *common.Response {
	ranks, err := s.rankMapper.FindByWebsiteId(websiteId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", ranks)
}

// GetRankDetail 获取榜单详情（含歌曲）
func (s *RankService) GetRankDetail(websiteId uint, rankName string) *common.Response {
	ranks, err := s.rankMapper.FindByWebsiteAndName(websiteId, rankName)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", ranks)
}