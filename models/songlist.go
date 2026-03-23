package models

import (
	"time"
)

type SongList struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:255" json:"title"`
	Pic         string    `gorm:"size:255" json:"pic"`
	Introduction string  `gorm:"type:text" json:"introduction"`
	Style       string    `gorm:"size:10" json:"style"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (SongList) TableName() string {
	return "song_list"
}
