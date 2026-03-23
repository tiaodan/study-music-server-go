package mapper

import (
	"study-music-server-go/models"
)

type BannerMapper struct{}

func NewBannerMapper() *BannerMapper {
	return &BannerMapper{}
}

func (*BannerMapper) FindAll() ([]models.Banner, error) {
	var banners []models.Banner
	err := DB.Order("create_time DESC").Find(&banners).Error
	return banners, err
}
