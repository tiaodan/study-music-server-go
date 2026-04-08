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

// 一键导入-单歌手-所有专辑请求
type ImportSingerAlbumsRequest struct {
	From string `json:"from"` // 歌手目录，如 C:\test\周杰伦（目录下有多个专辑子目录）
	To   string `json:"to"`   // 目标路径，如 \\100.86.118.11\hdd\周杰伦
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

// ===== S3 文件夹管理请求 =====

// CreateFolderRequest 创建文件夹
type CreateFolderRequest struct {
	Path string `json:"path"` // 文件夹路径，如 "music/album1/"
}

// DeleteFolderRequest 删除文件夹
type DeleteFolderRequest struct {
	Path      string `json:"path"`      // 文件夹路径
	Recursive bool   `json:"recursive"` // 是否递归删除（包含文件）
}

// DeleteObjectRequest 删除文件
type DeleteObjectRequest struct {
	Key string `json:"key"` // 文件完整路径（key）
}

// CopyObjectRequest 复制文件
type CopyObjectRequest struct {
	SourceKey string `json:"source_key"` // 源文件路径
	DestKey   string `json:"dest_key"`   // 目标文件路径
}

// BatchDeleteRequest 批量删除文件
type BatchDeleteRequest struct {
	Keys []string `json:"keys"` // 文件路径列表
}

// BatchCopyRequest 批量复制文件
type BatchCopyRequest struct {
	Items []CopyItem `json:"items"` // 复制项列表
}

// CopyItem 复制项
type CopyItem struct {
	SourceKey string `json:"source_key"` // 源文件路径
	DestKey   string `json:"dest_key"`   // 目标文件路径
}

// BatchGetInfoRequest 批量获取文件信息
type BatchGetInfoRequest struct {
	Keys []string `json:"keys"` // 文件路径列表
}

// ===== S3 响应结构 =====

// FolderInfo 文件夹信息
type FolderInfo struct {
	Path         string `json:"path"`          // 文件夹路径
	Name         string `json:"name"`          // 文件夹名称
	FileCount    int64  `json:"file_count"`    // 文件数量
	TotalSize    int64  `json:"total_size"`    // 总大小（字节）
	LastModified string `json:"last_modified"` // 最后修改时间
}

// ObjectInfo 文件对象信息
type ObjectInfo struct {
	Key          string `json:"key"`           // 完整路径
	Name         string `json:"name"`          // 文件名
	Size         int64  `json:"size"`          // 文件大小
	ContentType  string `json:"content_type"`  // MIME 类型
	LastModified string `json:"last_modified"` // 最后修改时间
	ETag         string `json:"etag"`          // ETag（MD5）
	URL          string `json:"url"`           // 访问 URL（临时签名 URL）
}

// ListFoldersResponse 文件夹列表响应
type ListFoldersResponse struct {
	Folders []FolderInfo `json:"folders"`
	Total   int          `json:"total"`
}

// ListObjectsResponse 文件列表响应
type ListObjectsResponse struct {
	Objects []ObjectInfo `json:"objects"`
	Total   int          `json:"total"`
}

// UploadFileResponse 上传文件响应
type UploadFileResponse struct {
	Key  string `json:"key"`  // 文件路径
	URL  string `json:"url"`  // 访问 URL
	Size int64  `json:"size"` // 文件大小
}

// BatchDeleteResponse 批量删除响应
type BatchDeleteResponse struct {
	Deleted []string `json:"deleted"` // 成功删除的文件
	Failed  []string `json:"failed"`  // 删除失败的文件
	Total   int      `json:"total"`   // 总数
}

// BatchCopyResponse 批量复制响应
type BatchCopyResponse struct {
	Copied []CopyResult `json:"copied"` // 成功复制的文件
	Failed []string     `json:"failed"` // 复制失败的文件
	Total  int          `json:"total"`  // 总数
}

// CopyResult 复制结果
type CopyResult struct {
	SourceKey string `json:"source_key"`
	DestKey   string `json:"dest_key"`
	URL       string `json:"url"`
}

// BatchGetInfoResponse 批量获取文件信息响应
type BatchGetInfoResponse struct {
	Objects []ObjectInfo `json:"objects"`
	Total   int          `json:"total"`
}
