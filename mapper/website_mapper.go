package mapper

import (
	"study-music-server-go/models"
)

type WebsiteMapper struct{}

func NewWebsiteMapper() *WebsiteMapper {
	return &WebsiteMapper{}
}

func (*WebsiteMapper) FindAll() ([]models.Website, error) {
	var websites []models.Website
	err := DB.Order("id").Find(&websites).Error
	return websites, err
}