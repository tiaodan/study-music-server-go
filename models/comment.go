package models

import (
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint      `json:"user_id"`
	SongId      uint      `json:"song_id"`
	SongListId  *uint     `json:"song_list_id"`
	Content     string    `gorm:"type:text" json:"content"`
	Type        *uint8    `gorm:"size:1" json:"type"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (Comment) TableName() string {
	return "comment"
}
