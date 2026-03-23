package mapper

import (
	"study-music-server-go/models"

	"gorm.io/gorm/clause"
)

type SongSingerMapper struct{}

func NewSongSingerMapper() *SongSingerMapper {
	return &SongSingerMapper{}
}

func (*SongSingerMapper) Add(songSinger *models.SongSinger) error {
	// UPSERT: 插入失败时更新（根据 song_id + singer_id 唯一索引）
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "song_id"}, {Name: "singer_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"update_time"}),
	}).Create(songSinger).Error
}

func (*SongSingerMapper) AddBatch(songSingers []models.SongSinger) error {
	return DB.Create(&songSingers).Error
}

func (*SongSingerMapper) FindBySongId(songId uint) ([]models.SongSinger, error) {
	var songSingers []models.SongSinger
	err := DB.Where("song_id = ?", songId).Find(&songSingers).Error
	return songSingers, err
}

func (*SongSingerMapper) FindBySingerId(singerId uint) ([]models.SongSinger, error) {
	var songSingers []models.SongSinger
	err := DB.Where("singer_id = ?", singerId).Find(&songSingers).Error
	return songSingers, err
}

func (*SongSingerMapper) DeleteBySongId(songId uint) error {
	return DB.Where("song_id = ?", songId).Delete(&models.SongSinger{}).Error
}

func (*SongSingerMapper) Delete(id uint) error {
	return DB.Delete(&models.SongSinger{}, id).Error
}
