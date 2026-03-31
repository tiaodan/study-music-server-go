package mapper

import (
	"study-music-server-go/models"
)

type SingerMapper struct{}

func NewSingerMapper() *SingerMapper {
	return &SingerMapper{}
}

func (*SingerMapper) Add(singer *models.Singer) error {
	return DB.Create(singer).Error
}

func (*SingerMapper) FindById(id uint) (*models.Singer, error) {
	var singer models.Singer
	err := DB.First(&singer, id).Error
	if err != nil {
		return nil, err
	}
	return &singer, nil
}

func (*SingerMapper) FindAll() ([]models.Singer, error) {
	var singers []models.Singer
	err := DB.Order("id").Find(&singers).Error
	return singers, err
}

// FindAllWithAlbums 只返回有专辑的歌手
func (*SingerMapper) FindAllWithAlbums() ([]models.Singer, error) {
	var singers []models.Singer
	err := DB.Distinct("singer.*").Table("singer").
		Joins("INNER JOIN album ON singer.id = album.singer_id").
		Order("singer.id").
		Find(&singers).Error
	return singers, err
}

func (*SingerMapper) FindByName(name string) ([]models.Singer, error) {
	var singers []models.Singer
	err := DB.Where("name LIKE ?", "%"+name+"%").Order("id").Find(&singers).Error
	return singers, err
}

// FindByNameWithAlbums 按名字搜索，只返回有专辑的歌手
func (*SingerMapper) FindByNameWithAlbums(name string) ([]models.Singer, error) {
	var singers []models.Singer
	err := DB.Distinct("singer.*").Table("singer").
		Joins("INNER JOIN album ON singer.id = album.singer_id").
		Where("singer.name LIKE ?", "%"+name+"%").
		Order("singer.id").
		Find(&singers).Error
	return singers, err
}

func (*SingerMapper) Update(singer *models.Singer) error {
	return DB.Save(singer).Error
}

func (*SingerMapper) Delete(id uint) error {
	return DB.Delete(&models.Singer{}, id).Error
}
