package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type RankListService struct {
	rankListMapper *mapper.RankListMapper
}

func NewRankListService() *RankListService {
	return &RankListService{
		rankListMapper: mapper.NewRankListMapper(),
	}
}

func (s *RankListService) AddRankList(req *models.RankListRequest) *common.Response {
	// Check if already rated
	existing, _ := s.rankListMapper.FindByUserIdAndSongListId(req.ConsumerId, req.SongListId)
	if existing != nil {
		// Update score
		existing.Score = &req.Score
		err := s.rankListMapper.Update(existing)
		if err != nil {
			return common.Error("评分失败")
		}
		return common.Success("评分更新成功")
	}

	rankList := &models.RankList{
		SongListId: req.SongListId,
		ConsumerId: req.ConsumerId,
		Score:      &req.Score,
	}
	err := s.rankListMapper.Add(rankList)
	if err != nil {
		return common.Error("评分失败")
	}
	return common.Success("评分成功")
}

func (s *RankListService) RankListOfSongListId(songListId uint) *common.Response {
	rankLists, err := s.rankListMapper.FindBySongListId(songListId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", rankLists)
}
