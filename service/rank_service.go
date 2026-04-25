package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

// parseInterval 解析时长字符串 "03:55" 返回秒数
func parseInterval(interval string) int {
	parts := strings.Split(interval, ":")
	if len(parts) != 2 {
		return 0
	}
	minutes, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	seconds, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil {
		return 0
	}
	return minutes * 60 + seconds
}

// removeBrackets 去掉歌名中的括号及括号内容
func removeBrackets(name string) string {
	re := strings.NewReplacer("（", "(", "）", ")")
	name = re.Replace(name)

	result := ""
	inBracket := false
	for _, ch := range name {
		if ch == '(' {
			inBracket = true
			continue
		}
		if ch == ')' {
			inBracket = false
			continue
		}
		if !inBracket {
			result += string(ch)
		}
	}

	return strings.TrimSpace(result)
}

// buildFileNameIndex 构建文件名索引
func buildFileNameIndex(files []utils.MusicFileInfo) map[string]string {
	index := make(map[string]string)
	for _, f := range files {
		singer := strings.ReplaceAll(f.Singer, " ", "")
		songName := strings.ReplaceAll(f.SongName, " ", "")
		key := singer + "-" + songName
		index[key] = f.OriginalName

		songNameNoBracket := removeBrackets(f.SongName)
		songNameNoBracket = strings.ReplaceAll(songNameNoBracket, " ", "")
		if songNameNoBracket != songName {
			altKey := singer + "-" + songNameNoBracket
			if _, exists := index[altKey]; !exists {
				index[altKey] = f.OriginalName
			}
		}
	}
	return index
}

// buildFilePathIndex 构建文件路径索引（用于查找 lrc 文件）
func buildFilePathIndex(files []utils.MusicFileInfo) map[string]string {
	index := make(map[string]string)
	for _, f := range files {
		singer := strings.ReplaceAll(f.Singer, " ", "")
		songName := strings.ReplaceAll(f.SongName, " ", "")
		key := singer + "-" + songName
		index[key] = f.Path

		songNameNoBracket := removeBrackets(f.SongName)
		songNameNoBracket = strings.ReplaceAll(songNameNoBracket, " ", "")
		if songNameNoBracket != songName {
			altKey := singer + "-" + songNameNoBracket
			if _, exists := index[altKey]; !exists {
				index[altKey] = f.Path
			}
		}
	}
	return index
}

// removeFileSpaces 去除文件名中的空格（mp3, wav, lrc 文件）
// 返回：重命名的文件数量，是否有错误
func removeFileSpaces(folderPath string) (int, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return 0, err
	}

	renamedCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		originalName := entry.Name()
		ext := strings.ToLower(filepath.Ext(originalName))

		// 只处理 mp3, wav, lrc 文件
		if ext != ".mp3" && ext != ".wav" && ext != ".lrc" {
			continue
		}

		// 检查文件名是否有空格
		if !strings.Contains(originalName, " ") {
			continue
		}

		// 去掉空格
		newName := strings.ReplaceAll(originalName, " ", "")
		oldPath := filepath.Join(folderPath, originalName)
		newPath := filepath.Join(folderPath, newName)

		// 重命名
		if err := os.Rename(oldPath, newPath); err != nil {
			log.Printf("重命名失败: %s -> %s, err: %v", oldPath, newPath, err)
			continue
		}

		renamedCount++
		log.Printf("重命名成功: %s -> %s", originalName, newName)
	}

	return renamedCount, nil
}

// ImportRank 导入排行榜数据
func (s *RankService) ImportRank(req *models.RankImportRequest) *common.Response {
	log.Printf("开始导入排行榜: websiteId=%d, rankName=%s, folderPath=%s, 歌曲数量=%d",
		req.WebsiteId, req.RankName, req.FolderPath, len(req.List))

	websiteName := s.getWebsiteName(req.WebsiteId)

	var importedSongs []map[string]interface{}
	var failed []string

	for _, item := range req.List {
		log.Printf("处理歌曲: %s - %s", item.Name, item.Singer)

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
					if album.ID == 0 {
						foundAlbum, _ = s.albumRankMapper.FindBySingerIdAndName(firstSingerId, albumName)
						if foundAlbum != nil {
							albumId = &foundAlbum.ID
						}
					} else {
						albumId = &album.ID
					}
				}
			}
		}

		// 3. 入库歌曲（nas_url_path 和 lyric 先为空）
		song := &models.SongRank{
			SingerId:       &firstSingerId,
			AlbumId:        albumId,
			Name:           strings.TrimSpace(item.Name),
			FullNameSinger: fullNameSinger,
			NasUrlPath:     "",
			Lyric:          "",
			SpiderUrl:      "",
			Duration:       parseInterval(item.Interval), // 从前端传入的时长
			Pic:            item.Img,                     // 从前端传入的封面
		}

		if err := s.songRankMapper.Add(song); err != nil {
			failed = append(failed, fmt.Sprintf("创建歌曲失败: %s - %v", item.Name, err))
			continue
		}

		// 4. 写入排行榜记录
		rank := &models.Rank{
			WebsiteId:   req.WebsiteId,
			Name:        req.RankName,
			SongRankId:  song.ID,
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

	// 第二阶段：匹配文件并读取歌词
	folderPath := strings.TrimSpace(req.FolderPath)
	var matchedCount int
	var lyricCount int

	if folderPath != "" {
		// 先处理文件名中的空格
		renamedCount, err := removeFileSpaces(folderPath)
		if err != nil {
			log.Printf("处理文件名空格失败: %v", err)
		} else if renamedCount > 0 {
			log.Printf("重命名 %d 个文件（去除空格）", renamedCount)
		}

		files, err := utils.GetMusicFiles(folderPath)
		if err != nil {
			log.Printf("读取目录失败: %v", err)
		} else {
			log.Printf("文件夹 %s 下有 %d 个音乐文件", folderPath, len(files))

			fileIndex := buildFileNameIndex(files)
			filePathIndex := buildFilePathIndex(files)
			log.Printf("文件名索引构建完成，共 %d 个", len(fileIndex))

			for i, songInfo := range importedSongs {
				fileKey := songInfo["file_key"].(string)
				fileName, found := fileIndex[fileKey]

				// 尝试去掉歌名括号匹配
				if !found {
					songName := songInfo["name"].(string)
					songNameNoBracket := removeBrackets(songName)
					singer := strings.ReplaceAll(songInfo["singer"].(string), " ", "")
					altKey := singer + "-" + strings.ReplaceAll(songNameNoBracket, " ", "")
					fileName, found = fileIndex[altKey]
				}

				// 尝试只用第一个歌手匹配
				if !found && strings.Contains(songInfo["singer"].(string), "、") {
					firstSinger := strings.Split(songInfo["singer"].(string), "、")[0]
					songName := songInfo["name"].(string)
					altKey := strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "") + "-" + strings.ReplaceAll(songName, " ", "")
					fileName, found = fileIndex[altKey]

					if !found {
						songNameNoBracket := removeBrackets(songName)
						altKey = strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "") + "-" + strings.ReplaceAll(songNameNoBracket, " ", "")
						fileName, found = fileIndex[altKey]
					}
				}

				if found {
					songId := songInfo["song_id"].(uint)
					nasUrlPath := fmt.Sprintf("rank/%s/%s/%s", websiteName, req.RankName, fileName)

					// 获取对应 mp3 文件的路径，查找 lrc 文件
					mp3Path := filePathIndex[fileKey]
					if mp3Path == "" {
						// 用其他 key 尝试获取路径
						songName := songInfo["name"].(string)
						songNameNoBracket := removeBrackets(songName)
						singer := strings.ReplaceAll(songInfo["singer"].(string), " ", "")
						mp3Path = filePathIndex[singer+"-"+strings.ReplaceAll(songNameNoBracket, " ", "")]

						if mp3Path == "" && strings.Contains(songInfo["singer"].(string), "、") {
							firstSinger := strings.Split(songInfo["singer"].(string), "、")[0]
							mp3Path = filePathIndex[strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "")+"-"+strings.ReplaceAll(songName, " ", "")]
						}
					}

					// 读取 lrc 文件
					var lyric string
					if mp3Path != "" {
						lrcPath := strings.TrimSuffix(mp3Path, ".mp3") + ".lrc"
						if lrcContent, err := readLrcFile(lrcPath); err == nil {
							lyric = lrcContent
							lyricCount++
							log.Printf("读取歌词成功: %s", lrcPath)
						} else {
							log.Printf("读取歌词失败: %s, err: %v", lrcPath, err)
						}
					}

					// 更新歌曲信息
					if err := s.songRankMapper.UpdateNasUrlPathAndLyric(songId, nasUrlPath, lyric); err != nil {
						log.Printf("更新歌曲信息失败: song_id=%d, err=%v", songId, err)
					} else {
						matchedCount++
						importedSongs[i]["nas_url_path"] = nasUrlPath
						importedSongs[i]["lyric"] = lyric
						log.Printf("匹配成功: %s -> %s", songInfo["name"], fileName)
					}
				}
			}
		}
	}

	log.Printf("第二阶段匹配完成: 匹配 %d 个文件, 读取 %d 个歌词", matchedCount, lyricCount)

	result := map[string]interface{}{
		"total":         len(req.List),
		"imported":      len(importedSongs),
		"failed":        len(failed),
		"matched_files": matchedCount,
		"lyric_files":   lyricCount,
		"results":       importedSongs,
		"failed_list":   failed,
	}

	log.Printf("排行榜导入完成: total=%d, imported=%d, failed=%d, matched=%d, lyric=%d",
		len(req.List), len(importedSongs), len(failed), matchedCount, lyricCount)

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

	// 填充 Singer 和 Album 计算字段
	for i := range ranks {
		if ranks[i].SongDetail != nil {
			// 填充歌手名
			if ranks[i].SongDetail.FullNameSinger != "" {
				ranks[i].SongDetail.Singer = ranks[i].SongDetail.FullNameSinger
			} else if ranks[i].SongDetail.SingerInfo != nil {
				ranks[i].SongDetail.Singer = ranks[i].SongDetail.SingerInfo.Name
			}

			// 填充专辑名
			if ranks[i].SongDetail.AlbumInfo != nil {
				ranks[i].SongDetail.Album = ranks[i].SongDetail.AlbumInfo.Name
			}
		}
	}

	return common.SuccessWithData("获取成功", ranks)
}

// CheckRank 校验排行榜数据（不入库，只对比文件）
func (s *RankService) CheckRank(req *models.RankCheckRequest) *common.Response {
	log.Printf("开始校验排行榜: websiteId=%d, rankName=%s, folderPath=%s, 歌曲数量=%d",
		req.WebsiteId, req.RankName, req.FolderPath, len(req.List))

	folderPath := strings.TrimSpace(req.FolderPath)
	if folderPath == "" {
		return common.Error("folderPath 不能为空")
	}

	files, err := utils.GetMusicFiles(folderPath)
	if err != nil {
		return common.Error(fmt.Sprintf("读取目录失败: %v", err))
	}
	log.Printf("文件夹 %s 下有 %d 个音乐文件", folderPath, len(files))

	fileIndex := buildFileNameIndex(files)
	matchedFileNames := make(map[string]bool)

	var missingFiles []models.MissingFileItem
	var matchedFiles []models.MatchedFileItem

	for _, item := range req.List {
		singer := strings.ReplaceAll(item.Singer, " ", "")
		songName := strings.ReplaceAll(item.Name, " ", "")
		key := singer + "-" + songName

		fileName, found := fileIndex[key]

		if !found {
			songNameNoBracket := removeBrackets(item.Name)
			altKey := singer + "-" + strings.ReplaceAll(songNameNoBracket, " ", "")
			fileName, found = fileIndex[altKey]
		}

		if !found && strings.Contains(item.Singer, "、") {
			firstSinger := strings.Split(item.Singer, "、")[0]
			altKey := strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "") + "-" + songName
			fileName, found = fileIndex[altKey]

			if !found {
				songNameNoBracket := removeBrackets(item.Name)
				altKey = strings.ReplaceAll(strings.TrimSpace(firstSinger), " ", "") + "-" + strings.ReplaceAll(songNameNoBracket, " ", "")
				fileName, found = fileIndex[altKey]
			}
		}

		if found {
			matchedFiles = append(matchedFiles, models.MatchedFileItem{
				Singer:   item.Singer,
				Name:     item.Name,
				FileName: fileName,
				Key:      key,
			})
			matchedFileNames[fileName] = true
		} else {
			missingFiles = append(missingFiles, models.MissingFileItem{
				Singer: item.Singer,
				Name:   item.Name,
				Key:    key,
			})
		}
	}

	var extraFiles []string
	for _, f := range files {
		if !matchedFileNames[f.OriginalName] {
			extraFiles = append(extraFiles, f.OriginalName)
		}
	}

	result := models.RankCheckResult{
		TotalFromFrontend: len(req.List),
		TotalFromFolder:   len(files),
		MissingFiles:      missingFiles,
		ExtraFiles:        extraFiles,
		MatchedFiles:      matchedFiles,
		CanImport:         len(missingFiles) == 0,
	}

	log.Printf("校验完成: 前端 %d 首, 文件夹 %d 个, 缺少 %d 个, 多余 %d 个, 匹配 %d 个",
		len(req.List), len(files), len(missingFiles), len(extraFiles), len(matchedFiles))

	return common.SuccessWithData("校验完成", result)
}