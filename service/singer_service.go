package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type SingerService struct {
	singerMapper *mapper.SingerMapper
	albumMapper  *mapper.AlbumMapper
	songMapper   *mapper.SongMapper
}

func NewSingerService() *SingerService {
	return &SingerService{
		singerMapper: mapper.NewSingerMapper(),
		albumMapper:  mapper.NewAlbumMapper(),
		songMapper:   mapper.NewSongMapper(),
	}
}

func (s *SingerService) AddSinger(req *models.SingerRequest) *common.Response {
	singer := &models.Singer{
		Name:        req.Name,
		Sex:         req.Sex,
		Pic:         req.Pic,
		Birth:       req.Birth,
		Location:    req.Location,
		Introduction: req.Introduction,
	}
	err := s.singerMapper.Add(singer)
	if err != nil {
		return common.Error("添加歌手失败")
	}
	return common.SuccessWithData("添加成功", singer)
}

func (s *SingerService) UpdateSinger(req *models.SingerRequest) *common.Response {
	singer, err := s.singerMapper.FindById(req.ID)
	if err != nil {
		return common.Error("歌手不存在")
	}
	singer.Name = req.Name
	singer.Sex = req.Sex
	singer.Pic = req.Pic
	singer.Birth = req.Birth
	singer.Location = req.Location
	singer.Introduction = req.Introduction
	err = s.singerMapper.Update(singer)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *SingerService) DeleteSinger(id uint) *common.Response {
	// 检查是否有歌曲关联（直接查 song 表的 singer_id）
	var count int64
	mapper.DB.Model(&models.Song{}).Where("singer_id = ?", id).Count(&count)
	if count > 0 {
		return common.Error("该歌手下有歌曲，无法删除")
	}
	err := s.singerMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}

func (s *SingerService) SingerOfId(id uint) *common.Response {
	singer, err := s.singerMapper.FindById(id)
	if err != nil {
		return common.Error("歌手不存在")
	}
	return common.SuccessWithData("获取成功", singer)
}

func (s *SingerService) SingerOfName(name string) *common.Response {
	singers, err := s.singerMapper.FindByNameWithAlbums(name)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", singers)
}

func (s *SingerService) AllSinger() *common.Response {
	singers, err := s.singerMapper.FindAllWithAlbums()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", singers)
}

// AlbumWithSongCount 专辑信息（含歌曲数量）
type AlbumWithSongCount struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	SingerId  uint   `json:"singer_id"`
	Pic       string `json:"pic"`
	SongCount int    `json:"song_count"`
}

// AlbumsBySingerId 查询歌手的所有专辑（含歌曲数量）
func (s *SingerService) AlbumsBySingerId(singerId uint) *common.Response {
	albums, err := s.albumMapper.FindBySingerId(singerId)
	if err != nil {
		return common.Error("获取专辑失败")
	}

	// 批量统计歌曲数量（避免N+1）
	var albumIds []uint
	for _, album := range albums {
		albumIds = append(albumIds, album.ID)
	}

	// 查询每个专辑的歌曲数量
	songCounts := make(map[uint]int)
	if len(albumIds) > 0 {
		type CountResult struct {
			AlbumId uint
			Count   int
		}
		var counts []CountResult
		mapper.DB.Table("song").
			Select("album_id, COUNT(*) as count").
			Where("album_id IN ?", albumIds).
			Group("album_id").
			Find(&counts)
		for _, c := range counts {
			songCounts[c.AlbumId] = c.Count
		}
	}

	// 组装结果
	result := make([]AlbumWithSongCount, len(albums))
	for i, album := range albums {
		result[i] = AlbumWithSongCount{
			ID:        album.ID,
			Name:      album.Name,
			SingerId:  album.SingerId,
			Pic:       album.Pic,
			SongCount: songCounts[album.ID],
		}
	}

	return common.SuccessWithData("获取成功", result)
}
