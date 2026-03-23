package models

import (
	"time"
)

type Admin struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string    `gorm:"size:50;uniqueIndex" json:"username"`
	Password    string    `gorm:"size:100" json:"-"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Admin) TableName() string {
	return "admin"
}
