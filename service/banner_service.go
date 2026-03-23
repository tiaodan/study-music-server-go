package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
)

type BannerService struct {
	bannerMapper *mapper.BannerMapper
}

func NewBannerService() *BannerService {
	return &BannerService{
		bannerMapper: mapper.NewBannerMapper(),
	}
}

func (s *BannerService) AllBanner() *common.Response {
	banners, err := s.bannerMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", banners)
}
