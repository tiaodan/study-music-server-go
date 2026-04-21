package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
)

type WebsiteService struct {
	websiteMapper *mapper.WebsiteMapper
}

func NewWebsiteService() *WebsiteService {
	return &WebsiteService{
		websiteMapper: mapper.NewWebsiteMapper(),
	}
}

func (s *WebsiteService) AllWebsite() *common.Response {
	websites, err := s.websiteMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", websites)
}