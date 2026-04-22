package models

import (
	"time"
)

// AlbumRank 排行榜专辑表
type AlbumRank struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"size:255;not null;uniqueIndex:idx_singer_name_rank" json:"name"` // 专辑名
	SingerId  uint      `gorm:"not null;uniqueIndex:idx_singer_name_rank" json:"singer_id"`     // 歌手ID
	Pic       string    `gorm:"size:255" json:"pic"`                                          // 封面
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`

	// 关联
	Singer *SingerRank `gorm:"foreignKey:SingerId" json:"singer,omitempty"`
}

func (AlbumRank) TableName() string {
	return "album_rank"
}