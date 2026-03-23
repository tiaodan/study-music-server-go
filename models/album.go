package models

import (
	"time"
)

type Album struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:100;not null" json:"name"`
	SingerId    uint       `gorm:"index;not null" json:"singer_id"`
	Pic         string     `gorm:"size:255" json:"pic"`
	Introduction string    `gorm:"size:255" json:"introduction"`
	CreateTime  time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time  `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Album) TableName() string {
	return "album"
}
