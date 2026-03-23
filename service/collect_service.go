package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type CollectService struct {
	collectMapper *mapper.CollectMapper
}

func NewCollectService() *CollectService {
	return &CollectService{
		collectMapper: mapper.NewCollectMapper(),
	}
}

func (s *CollectService) AddCollect(req *models.CollectRequest) *common.Response {
	// Check if already collected
	existing, _ := s.collectMapper.FindByUserIdAndSongId(req.UserId, req.SongId)
	if existing != nil {
		return common.Warning("已收藏")
	}

	collect := &models.Collect{
		UserId: req.UserId,
		SongId: req.SongId,
		Type:   req.Type,
	}
	err := s.collectMapper.Add(collect)
	if err != nil {
		return common.Error("收藏失败")
	}
	return common.Success("收藏成功")
}

func (s *CollectService) DeleteCollect(userId, songId uint) *common.Response {
	err := s.collectMapper.DeleteByUserIdAndSongId(userId, songId)
	if err != nil {
		return common.Error("取消收藏失败")
	}
	return common.Success("取消收藏成功")
}

func (s *CollectService) CollectOfUserId(userId uint) *common.Response {
	collects, err := s.collectMapper.FindByUserId(userId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", collects)
}
