package service

import (
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
)

type UserSupportService struct {
	userSupportMapper *mapper.UserSupportMapper
}

func NewUserSupportService() *UserSupportService {
	return &UserSupportService{
		userSupportMapper: mapper.NewUserSupportMapper(),
	}
}

func (s *UserSupportService) AddUserSupport(req *models.UserSupportRequest) *common.Response {
	// Check if already supported
	existing, _ := s.userSupportMapper.FindByUserIdAndCommentId(req.UserId, req.CommentId)
	if existing != nil {
		return common.Warning("已点赞")
	}

	userSupport := &models.UserSupport{
		UserId:    req.UserId,
		CommentId: req.CommentId,
		Type:      req.Type,
	}
	err := s.userSupportMapper.Add(userSupport)
	if err != nil {
		return common.Error("点赞失败")
	}
	return common.Success("点赞成功")
}

func (s *UserSupportService) DeleteUserSupport(userId, commentId uint) *common.Response {
	err := s.userSupportMapper.DeleteByUserIdAndCommentId(userId, commentId)
	if err != nil {
		return common.Error("取消点赞失败")
	}
	return common.Success("取消点赞成功")
}
