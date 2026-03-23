package models

import (
	"time"
)

type ListSong struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SongId      uint      `json:"song_id"`
	SongListId  uint      `json:"song_list_id"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
}

func (ListSong) TableName() string {
	return "list_song"
}
