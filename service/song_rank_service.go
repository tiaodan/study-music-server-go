package service

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
)

type SongRankService struct {
	songRankMapper *mapper.SongRankMapper
	deviceMapper   *mapper.DeviceMapper
}

func NewSongRankService() *SongRankService {
	return &SongRankService{
		songRankMapper: mapper.NewSongRankMapper(),
		deviceMapper:   mapper.NewDeviceMapper(),
	}
}

// SongOfId 获取歌曲详情
func (s *SongRankService) SongOfId(id uint) *common.Response {
	log.Printf("SongRankService.SongOfId 被调用, id=%d", id)

	song, err := s.songRankMapper.FindById(id)
	if err != nil {
		return common.Error("歌曲不存在")
	}

	// 填充歌手名
	if song.FullNameSinger != "" {
		song.Singer = song.FullNameSinger
	} else if song.SingerInfo != nil {
		song.Singer = song.SingerInfo.Name
	}

	// 填充专辑名
	if song.AlbumInfo != nil {
		song.Album = song.AlbumInfo.Name
	}

	// 返回播放 URL（通过音频流接口）
	song.Url = fmt.Sprintf("/stream/song-rank/%d", id)

	log.Printf("SongRank ID=%d, 歌名=%s, nas_url=%s, play_url=%s", id, song.Name, song.NasUrlPath, song.Url)

	return common.SuccessWithData("获取成功", song)
}

// GetFilePath 获取歌曲本地文件路径
func (s *SongRankService) GetFilePath(id uint) (string, error) {
	song, err := s.songRankMapper.FindById(id)
	if err != nil {
		return "", err
	}

	if song.NasUrlPath == "" {
		return "", os.ErrNotExist
	}

	nas, err := s.deviceMapper.FindByName("nas")
	if err != nil || nas == nil {
		return "", err
	}

	nasUrlPath := strings.TrimPrefix(song.NasUrlPath, "/")
	filePath := nas.UrlPrefix + "/" + nasUrlPath

	return filePath, nil
}

// StreamSong 读取歌曲文件并返回流
func (s *SongRankService) StreamSong(id uint) (*os.File, string, error) {
	filePath, err := s.GetFilePath(id)
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("打开文件失败: %s, err: %v", filePath, err)
		return nil, "", err
	}

	// 获取文件扩展名作为 content-type
	ext := strings.ToLower(filepath.Ext(filePath))
	var contentType string
	switch ext {
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	default:
		contentType = "application/octet-stream"
	}

	log.Printf("Streaming song id=%d, file=%s, contentType=%s", id, filePath, contentType)

	return file, contentType, nil
}

// GetLyric 获取歌词
func (s *SongRankService) GetLyric(id uint) (string, error) {
	song, err := s.songRankMapper.FindById(id)
	if err != nil {
		return "", err
	}
	return song.Lyric, nil
}