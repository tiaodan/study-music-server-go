package models

import (
	"time"
)

type Song struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AlbumId        *uint      `gorm:"index" json:"album_id"`
	Name           string     `gorm:"size:100;not null;uniqueIndex:idx_album_name" json:"name"`
	FullNameSinger string     `gorm:"size:255" json:"full_name_singer"`   // 多歌手时存储，单人则为空
	Introduction   string     `gorm:"size:255" json:"introduction"`
	CreateTime     time.Time  `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime     time.Time  `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
	Pic            string     `gorm:"size:255" json:"pic"`
	Lyric          string     `gorm:"type:text" json:"lyric"`
	NasUrlPath     string     `gorm:"size:255" json:"-"`        // NAS存储路径
	SpiderUrl      string     `gorm:"size:500" json:"-"`        // 爬取链接
	SpiderUrlHttps string     `gorm:"size:500" json:"-"`       // 带https的完整链接
	AwsUrl         string     `gorm:"size:500" json:"-"`       // AWS真实链接
	AwsUrlTemp     string     `gorm:"size:500" json:"-"`        // AWS临时链接
	IsHot          bool       `gorm:"default:false" json:"is_hot"` // 是否热门（用于vultr优先）

	// 计算字段，不存储到数据库
	Url       string `gorm:"-" json:"url"`       // 最优可用URL
	UrlSource string `gorm:"-" json:"url_source"` // URL来源：spider/vultr/aws/nas
	Singer    string `gorm:"-" json:"singer"`    // 歌手名（联表查询填充）
	Album     string `gorm:"-" json:"album"`     // 专辑名（联表查询填充）
}

func (Song) TableName() string {
	return "song"
}
