package mapper

import (
	"study-music-server-go/models"
)

type SingerRankMapper struct{}

func NewSingerRankMapper() *SingerRankMapper {
	return &SingerRankMapper{}
}

func (*SingerRankMapper) FindAll() ([]models.SingerRank, error) {
	var singers []models.SingerRank
	err := DB.Order("id").Find(&singers).Error
	return singers, err
}

func (*SingerRankMapper) FindById(id uint) (*models.SingerRank, error) {
	var singer models.SingerRank
	err := DB.First(&singer, id).Error
	if err != nil {
		return nil, err
	}
	return &singer, nil
}

func (*SingerRankMapper) FindByName(name string) ([]models.SingerRank, error) {
	var singers []models.SingerRank
	err := DB.Where("name LIKE ?", "%"+name+"%").Order("id").Find(&singers).Error
	return singers, err
}

func (*SingerRankMapper) Add(singer *models.SingerRank) error {
	return DB.Create(singer).Error
}

func (*SingerRankMapper) Update(singer *models.SingerRank) error {
	return DB.Save(singer).Error
}

func (*SingerRankMapper) Delete(id uint) error {
	return DB.Delete(&models.SingerRank{}, id).Error
}