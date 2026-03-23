package models

import "time"

type ConsumerRequest struct {
	ID          uint       `json:"id"`
	Username    string     `json:"username"`
	OldPassword string     `json:"old_password"`
	Password    string     `json:"password"`
	Sex         *uint8    `json:"sex"`
	PhoneNum    string     `json:"phone_num"`
	Email       string     `json:"email"`
	Birth       *time.Time `json:"birth"`
	Introduction string   `json:"introduction"`
	Location    string     `json:"location"`
	Avator      string     `json:"avator"`
	CreateTime  time.Time  `json:"create_time"`
}

type SingerRequest struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Sex         *uint8    `json:"sex"`
	Pic         string     `json:"pic"`
	Birth       *time.Time `json:"birth"`
	Location    string     `json:"location"`
	Introduction string   `json:"introduction"`
}

type SongRequest struct {
	ID             uint   `json:"id"`
	SingerId       uint   `json:"singer_id"`           // 保留用于兼容
	AlbumId        *uint  `json:"album_id"`            // 专辑ID
	Name           string `json:"name"`                // 歌曲名（不含歌手）
	Introduction   string `json:"introduction"`
	Pic            string `json:"pic"`
	Lyric          string `json:"lyric"`
	NasUrlPath     string `json:"nas_url_path"`        // NAS存储路径
	SpiderUrl      string `json:"spider_url"`          // 爬取链接（完整URL，带http头）
	SpiderUrlHttps string `json:"spider_url_https"`   // 带https的完整链接
	AwsUrl         string `json:"aws_url"`             // AWS真实链接（完整URL）
	AwsUrlTemp     string `json:"aws_url_temp"`        // AWS临时链接（完整URL）
	FullNameSinger string `json:"full_name_singer"`   // 多歌手时存储，单人则为空
}

// 名字格式化请求
type FormatNameRequest struct {
	Path string `json:"path"` // 歌手-专辑 路径
}

// 移动文件请求
// fromPath: 源目录路径，toPath: 目标根目录（会自动创建 歌手名/专辑名/ 子目录）
type MoveFileRequest struct {
	From string `json:"fromPath"` // 源目录路径，如 C:\test\周杰伦\哎呦，不错哦
	To   string `json:"toPath"`   // 目标根目录，如 D:\Music
}

// 规整进数据库请求
type ImportSongsRequest struct {
	Path string `json:"path"` // 要导入的文件夹路径
}

type SongListRequest struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Pic         string `json:"pic"`
	Introduction string `json:"introduction"`
	Style       string `json:"style"`
}

type CollectRequest struct {
	ID     uint `json:"id"`
	UserId uint `json:"user_id"`
	SongId uint `json:"song_id"`
	Type   *uint8 `json:"type"`
}

type CommentRequest struct {
	ID         uint   `json:"id"`
	UserId     uint   `json:"user_id"`
	SongId     uint   `json:"song_id"`
	SongListId *uint  `json:"song_list_id"`
	Content    string `json:"content"`
	Type       *uint8 `json:"type"`
}

type RankListRequest struct {
	ID         uint    `json:"id"`
	SongListId uint    `json:"song_list_id"`
	ConsumerId uint    `json:"consumer_id"`
	Score      float64 `json:"score"`
}

type AdminRequest struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ListSongRequest struct {
	ID         uint `json:"id"`
	SongId     uint `json:"song_id"`
	SongListId uint `json:"song_list_id"`
}

type UserSupportRequest struct {
	ID        uint   `json:"id"`
	UserId    uint   `json:"user_id"`
	CommentId uint   `json:"comment_id"`
	Type      *uint8 `json:"type"`
}

type ResetPasswordRequest struct {
	Email    string `json:"email"`
	Code     string `json:"code"`
	Password string `json:"password"`
}
