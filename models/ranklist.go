package models

import (
	"time"
)

type RankList struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SongListId  uint      `json:"song_list_id"`
	ConsumerId  uint      `json:"consumer_id"`
	Score       *float64  `gorm:"type:double" json:"score"`
	CreateTime  time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime  time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (RankList) TableName() string {
	return "rank_list"
}
