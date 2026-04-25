package models

import (
	"time"
)

// SongRank 排行榜歌曲表（结构和song一致，nas_url_path用/rank开头）
type SongRank struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SingerId       *uint     `gorm:"uniqueIndex:idx_singer_album_name_rank" json:"singer_id"`
	AlbumId        *uint     `gorm:"uniqueIndex:idx_singer_album_name_rank" json:"album_id"`
	Name           string    `gorm:"size:100;not null;uniqueIndex:idx_singer_album_name_rank" json:"name"`
	FullNameSinger string    `gorm:"size:255" json:"full_name_singer"`
	Introduction   string    `gorm:"size:255" json:"introduction"`
	Duration       int       `gorm:"default:0" json:"duration"`
	CreateTime     time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime     time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
	Pic            string    `gorm:"size:255" json:"pic"`
	Lyric          string    `gorm:"type:text" json:"lyric"`
	NasUrlPath     string    `gorm:"size:255" json:"nas_url_path"`
	SpiderUrl      string    `gorm:"size:500" json:"spider_url"`
	SpiderUrlHttps string    `gorm:"size:500" json:"spider_url_https"`
	AwsUrl         string    `gorm:"size:500" json:"aws_url"`
	AwsUrlTemp     string    `gorm:"size:500" json:"aws_url_temp"`
	IsHot          bool      `gorm:"default:false" json:"is_hot"`
	UploadAwsStatus int      `gorm:"default:0" json:"upload_aws_status"`

	// 关联
	SingerInfo *SingerRank `gorm:"foreignKey:SingerId" json:"singer_info,omitempty"`
	AlbumInfo  *AlbumRank  `gorm:"foreignKey:AlbumId" json:"album_info,omitempty"`

	// 计算字段
	Url       string `gorm:"-" json:"url"`
	UrlSource string `gorm:"-" json:"url_source"`
	Singer    string `gorm:"-" json:"singer"`
	Album     string `gorm:"-" json:"album"`
}

func (SongRank) TableName() string {
	return "song_rank"
}
