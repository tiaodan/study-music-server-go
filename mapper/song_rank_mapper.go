package mapper

import (
	"study-music-server-go/models"

	"gorm.io/gorm/clause"
)

type SongRankMapper struct{}

func NewSongRankMapper() *SongRankMapper {
	return &SongRankMapper{}
}

func (*SongRankMapper) Add(song *models.SongRank) error {
	// UPSERT: 插入失败时更新（根据唯一索引 singer_id + album_id + name）
	return DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "singer_id"}, {Name: "album_id"}, {Name: "name"}},
		DoUpdates: clause.AssignmentColumns([]string{"full_name_singer", "introduction", "pic", "nas_url_path", "update_time"}),
	}).Create(song).Error
}

func (*SongRankMapper) FindById(id uint) (*models.SongRank, error) {
	var song models.SongRank
	err := DB.First(&song, id).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (*SongRankMapper) FindBySingerIdAlbumIdName(singerId uint, albumId *uint, name string) (*models.SongRank, error) {
	var song models.SongRank
	query := DB.Where("singer_id = ? AND name = ?", singerId, name)
	if albumId != nil {
		query = query.Where("album_id = ?", *albumId)
	} else {
		query = query.Where("album_id IS NULL")
	}
	err := query.First(&song).Error
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func (*SongRankMapper) FindAll() ([]models.SongRank, error) {
	var songs []models.SongRank
	err := DB.Order("id").Find(&songs).Error
	return songs, err
}

func (*SongRankMapper) FindByAlbumId(albumId uint) ([]models.SongRank, error) {
	var songs []models.SongRank
	err := DB.Where("album_id = ?", albumId).Order("id").Find(&songs).Error
	return songs, err
}

func (*SongRankMapper) Update(song *models.SongRank) error {
	return DB.Save(song).Error
}

func (*SongRankMapper) UpdateNasUrlPath(id uint, nasUrlPath string) error {
	return DB.Model(&models.SongRank{}).Where("id = ?", id).Update("nas_url_path", nasUrlPath).Error
}

func (*SongRankMapper) UpdateNasUrlPathAndLyric(id uint, nasUrlPath string, lyric string) error {
	return DB.Model(&models.SongRank{}).Where("id = ?", id).Updates(map[string]interface{}{
		"nas_url_path": nasUrlPath,
		"lyric":        lyric,
	}).Error
}

func (*SongRankMapper) Delete(id uint) error {
	return DB.Delete(&models.SongRank{}, id).Error
}