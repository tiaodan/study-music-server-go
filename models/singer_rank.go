package models

import (
	"time"
)

// SingerRank 排行榜歌手表
type SingerRank struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null;uniqueIndex" json:"name"` // 歌手名
	Sex         int       `gorm:"default:0" json:"sex"`                       // 0-女 1-男 2-保密 3-组合
	Pic         string    `gorm:"size:255" json:"pic"`                        // 头像
	Birth       string    `gorm:"size:50" json:"birth"`                      // 生日
	Location    string    `gorm:"size:255" json:"location"`                  // 地区
	Introduction string    `gorm:"type:text" json:"introduction"`             // 简介
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (SingerRank) TableName() string {
	return "singer_rank"
}