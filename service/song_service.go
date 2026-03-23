package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
	"strings"
	"time"
)

type SongService struct {
	songMapper       *mapper.SongMapper
	songSingerMapper *mapper.SongSingerMapper
	singerMapper    *mapper.SingerMapper
	albumMapper      *mapper.AlbumMapper
	deviceMapper     *mapper.DeviceMapper
}

func NewSongService() *SongService {
	return &SongService{
		songMapper:       mapper.NewSongMapper(),
		songSingerMapper: mapper.NewSongSingerMapper(),
		singerMapper:    mapper.NewSingerMapper(),
		albumMapper:      mapper.NewAlbumMapper(),
		deviceMapper:     mapper.NewDeviceMapper(),
	}
}

// checkUrlAccessible 检查URL是否可访问
// 对HTTP URL使用HEAD请求，对本地/SMB路径使用文件存在检查
func (s *SongService) checkUrlAccessible(url string) bool {
	if url == "" {
		return false
	}

	// 检查是否为本地路径或SMB路径（以 \\ 或盘符 开头）
	isLocalPath := strings.HasPrefix(url, "\\\\") || strings.HasPrefix(url, "C:") ||
		strings.HasPrefix(url, "D:") || strings.HasPrefix(url, "E:") ||
		strings.HasPrefix(url, "F:") || strings.HasPrefix(url, "G:")

	if isLocalPath {
		// 本地路径/SMB路径：直接检查文件是否存在
		// UNC路径如 \\100.86.118.11\hdd 直接使用
		_, err := os.Stat(url)
		if err == nil {
			log.Printf("本地文件存在: %s", url)
			return true
		}
		log.Printf("本地文件不存在: %s, err: %v", url, err)
		return false
	}

	// HTTP URL：使用HEAD请求
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 400
}

// calculateBestUrlWithDetail 计算最优URL并返回详细信息
// 返回: url, detail, sourceType (spider/vultr/aws/nas)
func (s *SongService) calculateBestUrlWithDetail(song *models.Song) (string, string, string) {
	// 1. Spider（优先）
	if song.SpiderUrl != "" {
		detail := fmt.Sprintf("spider_url: %s", song.SpiderUrl)
		log.Printf("检查 spider_url: %s", song.SpiderUrl)
		if s.checkUrlAccessible(song.SpiderUrl) {
			log.Printf("最终使用 spider_url: %s", song.SpiderUrl)
			return song.SpiderUrl, detail, "spider"
		}
		log.Printf("spider_url 不可用")
	}
	if song.SpiderUrlHttps != "" {
		detail := fmt.Sprintf("spider_url_https: %s", song.SpiderUrlHttps)
		log.Printf("检查 spider_url_https: %s", song.SpiderUrlHttps)
		if s.checkUrlAccessible(song.SpiderUrlHttps) {
			log.Printf("最终使用 spider_url_https: %s", song.SpiderUrlHttps)
			return song.SpiderUrlHttps, detail, "spider"
		}
		log.Printf("spider_url_https 不可用")
	}

	// 2. vultr (cloud-service) 的 url_prefix + nas_url_path（仅当 is_hot=true 时）
	if song.NasUrlPath != "" && song.IsHot {
		cloud, err := s.deviceMapper.FindByName("vultr")
		if err == nil && cloud != nil {
			cloudUrl := cloud.UrlPrefix + "/" + song.NasUrlPath
			detail := fmt.Sprintf("vultr_url: %s", cloudUrl)
			log.Printf("检查 vultr_url: %s (vultr_prefix: %s, nas_url_path: %s)", cloudUrl, cloud.UrlPrefix, song.NasUrlPath)
			if s.checkUrlAccessible(cloudUrl) {
				log.Printf("最终使用 vultr_url: %s", cloudUrl)
				return cloudUrl, detail, "vultr"
			}
			log.Printf("vultr_url 不可用")
		} else {
			log.Printf("未找到 vultr 设备 或 is_hot=false")
		}
	}

	// 3. AWS URL
	if song.AwsUrl != "" {
		detail := fmt.Sprintf("aws_url: %s", song.AwsUrl)
		log.Printf("检查 aws_url: %s", song.AwsUrl)
		if s.checkUrlAccessible(song.AwsUrl) {
			log.Printf("最终使用 aws_url: %s", song.AwsUrl)
			return song.AwsUrl, detail, "aws"
		}
		log.Printf("aws_url 不可用")
	}
	if song.AwsUrlTemp != "" {
		detail := fmt.Sprintf("aws_url_temp: %s", song.AwsUrlTemp)
		log.Printf("检查 aws_url_temp: %s", song.AwsUrlTemp)
		if s.checkUrlAccessible(song.AwsUrlTemp) {
			log.Printf("最终使用 aws_url_temp: %s", song.AwsUrlTemp)
			return song.AwsUrlTemp, detail, "aws"
		}
		log.Printf("aws_url_temp 不可用")
	}

	// 4. nas 的 url_prefix + nas_url_path（最后备选）
	if song.NasUrlPath != "" {
		nas, err := s.deviceMapper.FindByName("nas")
		if err == nil && nas != nil {
			nasUrl := nas.UrlPrefix + "/" + song.NasUrlPath
			detail := fmt.Sprintf("nas_url: %s", nasUrl)
			log.Printf("检查 nas_url: %s (device_prefix: %s, nas_url_path: %s)", nasUrl, nas.UrlPrefix, song.NasUrlPath)
			if s.checkUrlAccessible(nasUrl) {
				log.Printf("最终使用 nas_url: %s", nasUrl)
				return nasUrl, detail, "nas"
			}
			log.Printf("nas_url 不可用")
		} else {
			log.Printf("未找到 nas 设备")
		}
	}

	return "", "无可用URL", ""
}

// calculateBestUrl 计算最优可用URL（自动验证可用性）
// 优先级：Spider > cloud-service > aws_url > nas
func (s *SongService) calculateBestUrl(song *models.Song) string {
	// 1. Spider（优先）
	if song.SpiderUrl != "" && s.checkUrlAccessible(song.SpiderUrl) {
		return song.SpiderUrl
	}
	if song.SpiderUrlHttps != "" && s.checkUrlAccessible(song.SpiderUrlHttps) {
		return song.SpiderUrlHttps
	}

	// 2. cloud-service 的 url_prefix + nas_url_path（仅当 is_hot=true 时）
	if song.NasUrlPath != "" && song.IsHot {
		cloud, err := s.deviceMapper.FindByName("vultr")
		if err == nil && cloud != nil {
			cloudUrl := cloud.UrlPrefix + "/" + song.NasUrlPath
			if s.checkUrlAccessible(cloudUrl) {
				return cloudUrl
			}
		}
	}

	// 3. AWS URL
	if song.AwsUrl != "" && s.checkUrlAccessible(song.AwsUrl) {
		return song.AwsUrl
	}
	if song.AwsUrlTemp != "" && s.checkUrlAccessible(song.AwsUrlTemp) {
		return song.AwsUrlTemp
	}

	// 4. nas 的 url_prefix + nas_url_path（最后备选）
	if song.NasUrlPath != "" {
		nas, err := s.deviceMapper.FindByName("nas")
		if err == nil && nas != nil {
			nasUrl := nas.UrlPrefix + "/" + song.NasUrlPath
			if s.checkUrlAccessible(nasUrl) {
				return nasUrl
			}
		}
	}

	return ""
}

func (s *SongService) AddSong(req *models.SongRequest) *common.Response {
	song := &models.Song{
		AlbumId:        req.AlbumId,
		Name:           req.Name,
		Introduction:   req.Introduction,
		Pic:            req.Pic,
		Lyric:          req.Lyric,
		NasUrlPath:     req.NasUrlPath,
		SpiderUrl:      req.SpiderUrl,
		SpiderUrlHttps: req.SpiderUrlHttps,
		AwsUrl:         req.AwsUrl,
		AwsUrlTemp:     req.AwsUrlTemp,
		FullNameSinger: req.FullNameSinger,
	}
	err := s.songMapper.Add(song)
	if err != nil {
		return common.Error("添加歌曲失败")
	}

	// 如果请求中指定了 SingerId，插入中间表关联
	if req.SingerId > 0 {
		songSinger := &models.SongSinger{
			SongId:   song.ID,
			SingerId: req.SingerId,
		}
		s.songSingerMapper.Add(songSinger)
	}

	return common.SuccessWithData("添加成功", song)
}

func (s *SongService) UpdateSong(req *models.SongRequest) *common.Response {
	song, err := s.songMapper.FindById(req.ID)
	if err != nil {
		return common.Error("歌曲不存在")
	}
	song.Name = req.Name
	song.Introduction = req.Introduction
	song.Pic = req.Pic
	song.Lyric = req.Lyric
	song.NasUrlPath = req.NasUrlPath
	song.SpiderUrl = req.SpiderUrl
	song.SpiderUrlHttps = req.SpiderUrlHttps
	song.AwsUrl = req.AwsUrl
	song.AwsUrlTemp = req.AwsUrlTemp
	song.FullNameSinger = req.FullNameSinger
	err = s.songMapper.Update(song)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *SongService) DeleteSong(id uint) *common.Response {
	err := s.songMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}

func (s *SongService) SongOfId(id uint) *common.Response {
	// 强制刷新日志
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println(">>> SongOfId 被调用, id:", id)

	song, err := s.songMapper.FindById(id)
	if err != nil {
		return common.Error("歌曲不存在")
	}

	// 打印歌曲信息用于调试（不打印歌词）
	log.Printf("===== Song ID: %d =====", id)
	log.Printf("SpiderUrl: %s", song.SpiderUrl)
	log.Printf("SpiderUrlHttps: %s", song.SpiderUrlHttps)
	log.Printf("NasUrlPath: %s", song.NasUrlPath)
	log.Printf("AwsUrl: %s", song.AwsUrl)
	log.Printf("AwsUrlTemp: %s", song.AwsUrlTemp)
	log.Printf("Lyric: (已跳过)")

	// 计算最优URL
	url, detail, source := s.calculateBestUrlWithDetail(song)
	song.Url = url
	song.UrlSource = source // 用于标识URL来源

	log.Printf("最终使用URL: %s", url)
	log.Printf("详细信息: %s", detail)
	log.Printf("URL来源: %s", source)

	return common.SuccessWithData("获取成功", song)
}

func (s *SongService) SongOfSingerId(singerId uint) *common.Response {
	// 通过中间表查询该歌手的所有歌曲
	songSingers, err := s.songSingerMapper.FindBySingerId(singerId)
	if err != nil {
		return common.Error("获取失败")
	}

	// 提取 song_id 列表
	var songIds []uint
	for _, ss := range songSingers {
		songIds = append(songIds, ss.SongId)
	}

	if len(songIds) == 0 {
		return common.SuccessWithData("获取成功", []models.Song{})
	}

	// 查询歌曲详情
	var songs []models.Song
	err = mapper.DB.Where("id IN ?", songIds).Find(&songs).Error
	if err != nil {
		return common.Error("获取失败")
	}

	// 填充歌手名、专辑名、url
	for i := range songs {
		songs[i].Url = fmt.Sprintf("/song/%d", songs[i].ID)

		// 填充专辑名
		if songs[i].AlbumId != nil {
			album, err := s.albumMapper.FindById(*songs[i].AlbumId)
			if err == nil && album != nil {
				songs[i].Album = album.Name
			}
		}
	}

	// 填充歌手名（当前歌手）
	singer, err := s.singerMapper.FindById(singerId)
	if err == nil && singer != nil {
		for i := range songs {
			songs[i].Singer = singer.Name
		}
	}

	return common.SuccessWithData("获取成功", songs)
}

func (s *SongService) SongOfName(name string) *common.Response {
	songs, err := s.songMapper.FindByName(name)
	if err != nil {
		return common.Error("获取失败")
	}
	// 计算每个歌曲的最优URL
	for i := range songs {
		songs[i].Url = s.calculateBestUrl(&songs[i])
	}
	return common.SuccessWithData("获取成功", songs)
}

func (s *SongService) AllSong() *common.Response {
	songs, err := s.songMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	// 计算每个歌曲的最优URL
	for i := range songs {
		songs[i].Url = s.calculateBestUrl(&songs[i])
	}
	return common.SuccessWithData("获取成功", songs)
}
