package mapper

import (
	"study-music-server-go/models"

	"gorm.io/gorm/clause"
)

type RankMapper struct{}

func NewRankMapper() *RankMapper {
	return &RankMapper{}
}

func (*RankMapper) FindAll() ([]models.Rank, error) {
	var ranks []models.Rank
	err := DB.Order("id").Find(&ranks).Error
	return ranks, err
}

func (*RankMapper) FindById(id uint) (*models.Rank, error) {
	var rank models.Rank
	err := DB.First(&rank, id).Error
	if err != nil {
		return nil, err
	}
	return &rank, nil
}

func (*RankMapper) FindByWebsiteAndName(websiteId uint, name string) ([]models.Rank, error) {
	var ranks []models.Rank
	err := DB.Preload("SongDetail").
		Preload("SongDetail.SingerInfo").
		Preload("SongDetail.AlbumInfo").
		Where("website_id = ? AND name = ?", websiteId, name).Order("id").Find(&ranks).Error
	return ranks, err
}

func (*RankMapper) FindByWebsiteId(websiteId uint) ([]models.Rank, error) {
	var ranks []models.Rank
	err := DB.Where("website_id = ?", websiteId).Order("id").Find(&ranks).Error
	return ranks, err
}

func (*RankMapper) Add(rank *models.Rank) error {
	// UPSERT: 插入失败时更新（根据唯一索引 website_id + name + song_rank_id）
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "website_id"}, {Name: "name"}, {Name: "song_rank_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"update_time"}),
	}).Create(rank).Error
}

func (*RankMapper) Update(rank *models.Rank) error {
	return DB.Save(rank).Error
}

func (*RankMapper) Delete(id uint) error {
	return DB.Delete(&models.Rank{}, id).Error
}

func (*RankMapper) DeleteByWebsiteAndName(websiteId uint, name string) error {
	return DB.Where("website_id = ? AND name = ?", websiteId, name).Delete(&models.Rank{}).Error
}