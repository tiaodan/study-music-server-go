package mapper

import (
	"study-music-server-go/models"
)

type ListSongMapper struct{}

func NewListSongMapper() *ListSongMapper {
	return &ListSongMapper{}
}

func (*ListSongMapper) Add(listSong *models.ListSong) error {
	return DB.Create(listSong).Error
}

func (*ListSongMapper) FindBySongListId(songListId uint) ([]models.ListSong, error) {
	var listSongs []models.ListSong
	err := DB.Where("song_list_id = ?", songListId).Find(&listSongs).Error
	return listSongs, err
}

func (*ListSongMapper) Delete(id uint) error {
	return DB.Delete(&models.ListSong{}, id).Error
}

func (*ListSongMapper) DeleteBySongListId(songListId uint) error {
	return DB.Where("song_list_id = ?", songListId).Delete(&models.ListSong{}).Error
}
