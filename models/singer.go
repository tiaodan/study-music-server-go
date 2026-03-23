package models

import (
	"time"
)

type Singer struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:45" json:"name"`
	Sex         *uint8    `gorm:"size:1" json:"sex"`
	Pic         string    `gorm:"size:255" json:"pic"`
	Birth       *time.Time `json:"birth"`
	Location    string    `gorm:"size:45" json:"location"`
	Introduction string  `gorm:"size:255" json:"introduction"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Singer) TableName() string {
	return "singer"
}
