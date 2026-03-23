package models

import (
	"time"
)

type Collect struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserId      uint      `json:"user_id"`
	SongId      uint      `json:"song_id"`
	Type        *uint8   `gorm:"size:1" json:"type"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (Collect) TableName() string {
	return "collect"
}
