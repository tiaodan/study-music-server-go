package models

import (
	"time"
)

// Website 音乐网站表
type Website struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex:idx_name_type" json:"name"`   // 网站名称：QQ音乐、酷狗音乐等
	Type      string    `gorm:"size:20;not null;uniqueIndex:idx_name_type" json:"type"`   // 类型：music
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Website) TableName() string {
	return "website"
}