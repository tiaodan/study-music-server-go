package models

import (
	"time"
)

type UserSupport struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint      `json:"user_id"`
	CommentId   uint      `json:"comment_id"`
	Type        *uint8    `gorm:"size:1" json:"type"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (UserSupport) TableName() string {
	return "user_support"
}
