package models

import (
	"time"
)

type Consumer struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"size:50" json:"username"`
	Password    string    `gorm:"size:100" json:"-"`
	Sex         *uint8    `gorm:"size:1" json:"sex"`
	PhoneNum    string    `gorm:"size:20" json:"phone_num"`
	Email       string    `gorm:"size:50" json:"email"`
	Birth       *time.Time `json:"birth"`
	Introduction string  `gorm:"size:255" json:"introduction"`
	Location    string    `gorm:"size:50" json:"location"`
	Avator      string    `gorm:"size:255" json:"avator"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Consumer) TableName() string {
	return "consumer"
}
