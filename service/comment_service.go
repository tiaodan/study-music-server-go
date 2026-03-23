package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type CommentService struct {
	commentMapper *mapper.CommentMapper
}

func NewCommentService() *CommentService {
	return &CommentService{
		commentMapper: mapper.NewCommentMapper(),
	}
}

func (s *CommentService) AddComment(req *models.CommentRequest) *common.Response {
	comment := &models.Comment{
		UserId:     req.UserId,
		SongId:     req.SongId,
		SongListId: req.SongListId,
		Content:    req.Content,
		Type:       req.Type,
	}
	err := s.commentMapper.Add(comment)
	if err != nil {
		return common.Error("评论失败")
	}
	return common.SuccessWithData("评论成功", comment)
}

func (s *CommentService) DeleteComment(id uint) *common.Response {
	err := s.commentMapper.Delete(id)
	if err != nil {
		return common.Error("删除失败")
	}
	return common.Success("删除成功")
}

func (s *CommentService) CommentOfSongId(songId uint) *common.Response {
	comments, err := s.commentMapper.FindBySongId(songId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", comments)
}

func (s *CommentService) CommentOfSongListId(songListId uint) *common.Response {
	comments, err := s.commentMapper.FindBySongListId(songListId)
	if err != nil {
		return common.Error("获取失败")
	}
	return common.SuccessWithData("获取成功", comments)
}
