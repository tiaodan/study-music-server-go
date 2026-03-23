package models

import (
	"time"
)

// SongSinger 中间表 - 歌曲与歌手多对多关系
type SongSinger struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SongId     uint      `gorm:"uniqueIndex:idx_song_singer;not null" json:"song_id"`
	SingerId   uint      `gorm:"uniqueIndex:idx_song_singer;not null" json:"singer_id"`
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (SongSinger) TableName() string {
	return "song_singer"
}
