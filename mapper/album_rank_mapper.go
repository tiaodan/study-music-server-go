package mapper

import (
	"study-music-server-go/models"

	"gorm.io/gorm/clause"
)

type AlbumRankMapper struct{}

func NewAlbumRankMapper() *AlbumRankMapper {
	return &AlbumRankMapper{}
}

func (*AlbumRankMapper) FindAll() ([]models.AlbumRank, error) {
	var albums []models.AlbumRank
	err := DB.Order("id").Find(&albums).Error
	return albums, err
}

func (*AlbumRankMapper) FindById(id uint) (*models.AlbumRank, error) {
	var album models.AlbumRank
	err := DB.First(&album, id).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (*AlbumRankMapper) FindBySingerId(singerId uint) ([]models.AlbumRank, error) {
	var albums []models.AlbumRank
	err := DB.Where("singer_id = ?", singerId).Order("id").Find(&albums).Error
	return albums, err
}

func (*AlbumRankMapper) FindBySingerIdAndName(singerId uint, name string) (*models.AlbumRank, error) {
	var album models.AlbumRank
	err := DB.Where("singer_id = ? AND name = ?", singerId, name).First(&album).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (*AlbumRankMapper) Add(album *models.AlbumRank) error {
	// UPSERT: 插入失败时更新（根据唯一索引 singer_id + name）
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "singer_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"pic", "update_time"}),
	}).Create(album).Error
}

func (*AlbumRankMapper) Update(album *models.AlbumRank) error {
	return DB.Save(album).Error
}

func (*AlbumRankMapper) Delete(id uint) error {
	return DB.Delete(&models.AlbumRank{}, id).Error
}