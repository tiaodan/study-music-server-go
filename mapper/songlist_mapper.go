package mapper

import (
	"study-music-server-go/models"
)

type SongListMapper struct{}

func NewSongListMapper() *SongListMapper {
	return &SongListMapper{}
}

func (*SongListMapper) Add(songList *models.SongList) error {
	return DB.Create(songList).Error
}

func (*SongListMapper) FindById(id uint) (*models.SongList, error) {
	var songList models.SongList
	err := DB.First(&songList, id).Error
	if err != nil {
		return nil, err
	}
	return &songList, nil
}

func (*SongListMapper) FindAll() ([]models.SongList, error) {
	var songLists []models.SongList
	err := DB.Find(&songLists).Error
	return songLists, err
}

func (*SongListMapper) FindByTitle(title string) ([]models.SongList, error) {
	var songLists []models.SongList
	err := DB.Where("title LIKE ?", "%"+title+"%").Find(&songLists).Error
	return songLists, err
}

func (*SongListMapper) Update(songList *models.SongList) error {
	return DB.Save(songList).Error
}

func (*SongListMapper) Delete(id uint) error {
	return DB.Delete(&models.SongList{}, id).Error
}
