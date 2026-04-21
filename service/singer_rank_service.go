package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type SingerRankService struct {
	singerRankMapper *mapper.SingerRankMapper
	albumRankMapper  *mapper.AlbumRankMapper
}

func NewSingerRankService() *SingerRankService {
	return &SingerRankService{
		singerRankMapper: mapper.NewSingerRankMapper(),
		albumRankMapper:  mapper.NewAlbumRankMapper(),
	}
}

func (s *SingerRankService) AllSinger() *common.Response {
	singers, err := s.singerRankMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", singers)
}

func (s *SingerRankService) SingerOfId(id uint) *common.Response {
	singer, err := s.singerRankMapper.FindById(id)
	if err != nil {
		return common.Error("歌手不存在")
	}
	return common.SuccessWithData("获取成功", singer)
}

func (s *SingerRankService) SingerOfName(name string) *common.Response {
	singers, err := s.singerRankMapper.FindByName(name)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", singers)
}

func (s *SingerRankService) AddSinger(req *models.SingerRankRequest) *common.Response {
	singer := &models.SingerRank{
		Name:         req.Name,
		Sex:          req.Sex,
		Pic:          req.Pic,
		Birth:        req.Birth,
		Location:     req.Location,
		Introduction: req.Introduction,
	}
	err := s.singerRankMapper.Add(singer)
	if err != nil {
		return common.Error("添加失败")
	}
	return common.SuccessWithData("添加成功", singer)
}

func (s *SingerRankService) UpdateSinger(req *models.SingerRankRequest) *common.Response {
	singer, err := s.singerRankMapper.FindById(req.ID)
	if err != nil {
		return common.Error("歌手不存在")
	}
	singer.Name = req.Name
	singer.Sex = req.Sex
	singer.Pic = req.Pic
	singer.Birth = req.Birth
	singer.Location = req.Location
	singer.Introduction = req.Introduction
	err = s.singerRankMapper.Update(singer)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *SingerRankService) DeleteSinger(id uint) *common.Response {
	err := s.singerRankMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}