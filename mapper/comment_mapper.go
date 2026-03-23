package mapper

import (
	"study-music-server-go/models"
)

type CommentMapper struct{}

func NewCommentMapper() *CommentMapper {
	return &CommentMapper{}
}

func (*CommentMapper) Add(comment *models.Comment) error {
	return DB.Create(comment).Error
}

func (*CommentMapper) FindById(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := DB.First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (*CommentMapper) FindBySongId(songId uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := DB.Where("song_id = ?", songId).Order("create_time DESC").Find(&comments).Error
	return comments, err
}

func (*CommentMapper) FindBySongListId(songListId uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := DB.Where("song_list_id = ?", songListId).Order("create_time DESC").Find(&comments).Error
	return comments, err
}

func (*CommentMapper) Delete(id uint) error {
	return DB.Delete(&models.Comment{}, id).Error
}
