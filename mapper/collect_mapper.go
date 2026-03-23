package mapper

import (
	"study-music-server-go/models"
)

type CollectMapper struct{}

func NewCollectMapper() *CollectMapper {
	return &CollectMapper{}
}

func (*CollectMapper) Add(collect *models.Collect) error {
	return DB.Create(collect).Error
}

func (*CollectMapper) FindById(id uint) (*models.Collect, error) {
	var collect models.Collect
	err := DB.First(&collect, id).Error
	if err != nil {
		return nil, err
	}
	return &collect, nil
}

func (*CollectMapper) FindByUserId(userId uint) ([]models.Collect, error) {
	var collects []models.Collect
	err := DB.Where("user_id = ?", userId).Find(&collects).Error
	return collects, err
}

func (*CollectMapper) FindByUserIdAndSongId(userId, songId uint) (*models.Collect, error) {
	var collect models.Collect
	err := DB.Where("user_id = ? AND song_id = ?", userId, songId).First(&collect).Error
	if err != nil {
		return nil, err
	}
	return &collect, nil
}

func (*CollectMapper) Delete(id uint) error {
	return DB.Delete(&models.Collect{}, id).Error
}

func (*CollectMapper) DeleteByUserIdAndSongId(userId, songId uint) error {
	return DB.Where("user_id = ? AND song_id = ?", userId, songId).Delete(&models.Collect{}).Error
}
