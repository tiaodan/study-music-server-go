package mapper

import (
	"study-music-server-go/models"
)

type RankListMapper struct{}

func NewRankListMapper() *RankListMapper {
	return &RankListMapper{}
}

func (*RankListMapper) Add(rankList *models.RankList) error {
	return DB.Create(rankList).Error
}

func (*RankListMapper) FindBySongListId(songListId uint) ([]models.RankList, error) {
	var rankLists []models.RankList
	err := DB.Where("song_list_id = ?", songListId).Order("score DESC").Find(&rankLists).Error
	return rankLists, err
}

func (*RankListMapper) FindByUserIdAndSongListId(userId, songListId uint) (*models.RankList, error) {
	var rankList models.RankList
	err := DB.Where("consumer_id = ? AND song_list_id = ?", userId, songListId).First(&rankList).Error
	if err != nil {
		return nil, err
	}
	return &rankList, nil
}

func (*RankListMapper) Update(rankList *models.RankList) error {
	return DB.Save(rankList).Error
}
