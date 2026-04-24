package models

import (
	"time"
)

// Rank 排行榜表
type Rank struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	WebsiteId uint      `gorm:"not null;uniqueIndex:idx_rank_unique" json:"website_id"`
	Name      string    `gorm:"size:50;not null;uniqueIndex:idx_rank_unique" json:"name"`
	SongRankId uint     `gorm:"not null;uniqueIndex:idx_rank_unique;column:song_rank_id" json:"song_rank_id"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`

	// 关联
	Website   *Website    `gorm:"foreignKey:WebsiteId" json:"website,omitempty"`
	SongDetail *SongRank  `gorm:"foreignKey:SongRankId" json:"song_detail,omitempty"`
}

func (Rank) TableName() string {
	return "rank"
}