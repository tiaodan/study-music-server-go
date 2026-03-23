package models

import "time"

// Device 存储设备
type Device struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string    `gorm:"size:100;not null" json:"name"`          // 设备名称，如 "我的NAS"、"Vultr云盘"
	Type       string    `gorm:"size:20;not null" json:"type"`          // 设备类型：nas / cloud-service
	UrlPrefix  string    `gorm:"size:255;not null" json:"url_prefix"`   // 路径前缀
	IsDefault  bool      `gorm:"default:false" json:"is_default"`       // 是否为默认存储
	CreateTime time.Time `gorm:"autoCreateTime" json:"create_time"`
	UpdateTime time.Time `gorm:"autoCreateTime;autoUpdateTime" json:"update_time"`
}

func (Device) TableName() string {
	return "device"
}
