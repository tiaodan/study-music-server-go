package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type AlbumRankService struct {
	albumRankMapper *mapper.AlbumRankMapper
}

func NewAlbumRankService() *AlbumRankService {
	return &AlbumRankService{
		albumRankMapper: mapper.NewAlbumRankMapper(),
	}
}

func (s *AlbumRankService) AllAlbum() *common.Response {
	albums, err := s.albumRankMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", albums)
}

func (s *AlbumRankService) AlbumOfId(id uint) *common.Response {
	album, err := s.albumRankMapper.FindById(id)
	if err != nil {
		return common.Error("专辑不存在")
	}
	return common.SuccessWithData("获取成功", album)
}

func (s *AlbumRankService) AlbumsOfSingerId(singerId uint) *common.Response {
	albums, err := s.albumRankMapper.FindBySingerId(singerId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", albums)
}

func (s *AlbumRankService) AddAlbum(req *models.AlbumRankRequest) *common.Response {
	album := &models.AlbumRank{
		Name:     req.Name,
		SingerId: req.SingerId,
		Pic:      req.Pic,
	}
	err := s.albumRankMapper.Add(album)
	if err != nil {
		return common.Error("添加失败")
	}
	return common.SuccessWithData("添加成功", album)
}

func (s *AlbumRankService) UpdateAlbum(req *models.AlbumRankRequest) *common.Response {
	album, err := s.albumRankMapper.FindById(req.ID)
	if err != nil {
		return common.Error("专辑不存在")
	}
	album.Name = req.Name
	album.SingerId = req.SingerId
	album.Pic = req.Pic
	err = s.albumRankMapper.Update(album)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *AlbumRankService) DeleteAlbum(id uint) *common.Response {
	err := s.albumRankMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}