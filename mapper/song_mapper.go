package mapper

import (
	"study-music-server-go/models"

	"gorm.io/gorm/clause"
)

type SongMapper struct{}

func NewSongMapper() *SongMapper {
	return &SongMapper{}
}

func (*SongMapper) Add(song *models.Song) error {
	// UPSERT: 插入失败时更新（根据唯一索引）
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"album_id", "nas_url_path", "full_name_singer", "introduction", "pic", "update_time"}),
	}).Create(song).Error
}

func (*SongMapper) FindById(id uint) (*models.Song, error) {
	var song models.Song
	err := DB.First(&song, id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (*SongMapper) FindAll() ([]models.Song, error) {
	var songs []models.Song
	err := DB.Find(&songs).Error
	return songs, err
}

func (*SongMapper) FindByAlbumId(albumId uint) ([]models.Song, error) {
	var songs []models.Song
	err := DB.Where("album_id = ?", albumId).Find(&songs).Error
	return songs, err
}

func (*SongMapper) FindByName(name string) ([]models.Song, error) {
	var songs []models.Song
	err := DB.Where("name LIKE ?", "%"+name+"%").Find(&songs).Error
	return songs, err
}

func (*SongMapper) Update(song *models.Song) error {
	return DB.Save(song).Error
}

func (*SongMapper) Delete(id uint) error {
	return DB.Delete(&models.Song{}, id).Error
}
