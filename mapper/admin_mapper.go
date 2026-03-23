package mapper

import (
	"study-music-server-go/models"
)

type AdminMapper struct{}

func NewAdminMapper() *AdminMapper {
	return &AdminMapper{}
}

func (*AdminMapper) FindByUsername(username string) (*models.Admin, error) {
	var admin models.Admin
	err := DB.Where("username = ?", username).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

func (*AdminMapper) FindById(id uint) (*models.Admin, error) {
	var admin models.Admin
	err := DB.First(&admin, id).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
