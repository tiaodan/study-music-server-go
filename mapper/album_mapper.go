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
	err := DB.Order("id desc").Find(&albums).Error
	return albums, err
}

func (*AlbumMapper) Update(album *models.Album) error {
	return DB.Save(album).Error
}

func (*AlbumMapper) Delete(id uint) error {
	return DB.Delete(&models.Album{}, id).Error
}

// FindBySingerId 查询歌手的所有专辑（按ID倒序）
func (*AlbumMapper) FindBySingerId(singerId uint) ([]models.Album, error) {
	var albums []models.Album
	err := DB.Where("singer_id = ?", singerId).Order("id desc").Find(&albums).Error
	return albums, err
}
