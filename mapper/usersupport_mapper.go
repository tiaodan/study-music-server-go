package mapper

import (
	"study-music-server-go/models"
)

type UserSupportMapper struct{}

func NewUserSupportMapper() *UserSupportMapper {
	return &UserSupportMapper{}
}

func (*UserSupportMapper) Add(userSupport *models.UserSupport) error {
	return DB.Create(userSupport).Error
}

func (*UserSupportMapper) FindByUserIdAndCommentId(userId, commentId uint) (*models.UserSupport, error) {
	var userSupport models.UserSupport
	err := DB.Where("user_id = ? AND comment_id = ?", userId, commentId).First(&userSupport).Error
	if err != nil {
		return nil, err
	}
	return &userSupport, nil
}

func (*UserSupportMapper) Delete(id uint) error {
	return DB.Delete(&models.UserSupport{}, id).Error
}

func (*UserSupportMapper) DeleteByUserIdAndCommentId(userId, commentId uint) error {
	return DB.Where("user_id = ? AND comment_id = ?", userId, commentId).Delete(&models.UserSupport{}).Error
}
