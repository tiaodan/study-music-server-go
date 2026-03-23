package mapper

import (
	"study-music-server-go/models"
)

type AlbumMapper struct{}

func NewAlbumMapper() *AlbumMapper {
	return &AlbumMapper{}
}

func (*AlbumMapper) Add(album *models.Album) error {
	return DB.Create(album).Error
}

func (*AlbumMapper) FindById(id uint) (*models.Album, error) {
	var album models.Album
	err := DB.First(&album, id).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (*AlbumMapper) FindByNameAndSingerId(name string, singerId uint) (*models.Album, error) {
	var album models.Album
	err := DB.Where("name = ? AND singer_id = ?", name, singerId).First(&album).Error
	if err != nil {
		return nil, err
	}
	return &album, nil
}

func (*AlbumMapper) FindAll() ([]models.Album, error) {
	var albums []models.Album
	err := DB.Find(&albums).Error
	return albums, err
}

func (*AlbumMapper) Update(album *models.Album) error {
	return DB.Save(album).Error
}

func (*AlbumMapper) Delete(id uint) error {
	return DB.Delete(&models.Album{}, id).Error
}
