package models

import (
	"time"
)

// Rank 排行榜表
type Rank struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	WebsiteId  uint      `gorm:"not null;uniqueIndex:idx_website_name_song" json:"website_id"` // 网站ID（外键）
	Name       string    `gorm:"size:50;not null;uniqueIndex:idx_website_name_song" json:"name"` // 榜单名：top500、热歌榜等
	SongId     uint      `gorm:"not null;uniqueIndex:idx_website_name_song" json:"song_id"`     // 歌曲ID（song_rank表外键）
	AlbumId    *uint     `gorm:"index" json:"album_id"`                                        // 专辑ID（album_rank表外键，冗余存储）
	Album      string    `gorm:"size:100" json:"album"`                                         // 专辑名
	Singer     string    `gorm:"size:255" json:"singer"`                                       // 歌手名（单人用album歌手，多人用full_name_singer）
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`

	// 关联
	Website *Website  `gorm:"foreignKey:WebsiteId" json:"website,omitempty"`
	Song    *SongRank `gorm:"foreignKey:SongId" json:"song_detail,omitempty"`
}

func (Rank) TableName() string {
	return "rank"
}