package models

import (
	"time"
)

type Banner struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string    `gorm:"size:50" json:"title"`
	Pic         string    `gorm:"size:255" json:"pic"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Banner) TableName() string {
	return "banner"
}
