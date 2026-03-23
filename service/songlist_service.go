package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type SongListService struct {
	songListMapper *mapper.SongListMapper
}

func NewSongListService() *SongListService {
	return &SongListService{
		songListMapper: mapper.NewSongListMapper(),
	}
}

func (s *SongListService) AddSongList(req *models.SongListRequest) *common.Response {
	songList := &models.SongList{
		Title:       req.Title,
		Pic:         req.Pic,
		Introduction: req.Introduction,
		Style:       req.Style,
	}
	err := s.songListMapper.Add(songList)
	if err != nil {
		return common.Error("添加歌单失败")
	}
	return common.SuccessWithData("添加成功", songList)
}

func (s *SongListService) UpdateSongList(req *models.SongListRequest) *common.Response {
	songList, err := s.songListMapper.FindById(req.ID)
	if err != nil {
		return common.Error("歌单不存在")
	}
	songList.Title = req.Title
	songList.Pic = req.Pic
	songList.Introduction = req.Introduction
	songList.Style = req.Style
	err = s.songListMapper.Update(songList)
	if err != nil {
		return common.Error("更新失败")
	}
	return common.Success("更新成功")
}

func (s *SongListService) DeleteSongList(id uint) *common.Response {
	err := s.songListMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}

func (s *SongListService) SongListOfId(id uint) *common.Response {
	songList, err := s.songListMapper.FindById(id)
	if err != nil {
		return common.Error("歌单不存在")
	}
	return common.SuccessWithData("获取成功", songList)
}

func (s *SongListService) SongListOfTitle(title string) *common.Response {
	songLists, err := s.songListMapper.FindByTitle(title)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", songLists)
}

func (s *SongListService) AllSongList() *common.Response {
	songLists, err := s.songListMapper.FindAll()
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", songLists)
}
