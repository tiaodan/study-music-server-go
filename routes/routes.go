package routes

import (
	"study-music-server-go/controller"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// User routes
	consumerCtrl := controller.NewConsumerController()
	r.POST("/user/add", consumerCtrl.AddUser)
	r.POST("/user/login/status", consumerCtrl.LoginStatus)
	r.POST("/user/email/status", consumerCtrl.LoginEmailStatus)
	r.GET("/user", consumerCtrl.AllUser)
	r.GET("/user/detail", consumerCtrl.UserOfId)
	r.GET("/user/delete", consumerCtrl.DeleteUser)
	r.POST("/user/update", consumerCtrl.UpdateUserMsg)
	r.POST("/user/updatePassword", consumerCtrl.UpdatePassword)
	r.POST("/user/avatar/update", consumerCtrl.UpdateUserAvatar)

	// Singer routes
	singerCtrl := controller.NewSingerController()
	r.POST("/singer/add", singerCtrl.AddSinger)
	r.POST("/singer/update", singerCtrl.UpdateSinger)
	r.GET("/singer/delete", singerCtrl.DeleteSinger)
	r.GET("/singer/detail", singerCtrl.SingerOfId)
	r.GET("/singer/name/detail", singerCtrl.SingerOfName)
	r.GET("/singer", singerCtrl.AllSinger)
	r.GET("/singer/jay", singerCtrl.SingerJay) // 临时：只返回周杰伦
	r.GET("/singer/albums", singerCtrl.AlbumsOfSingerId) // 新增：查歌手专辑列表

	// SongRank routes - 排行榜歌曲（直接返回音频流，和歌手歌曲逻辑一致）
	songRankCtrl := controller.NewSongRankController()
	r.GET("/song-rank/:id", songRankCtrl.SongOfId)

	// Song routes
	songCtrl := controller.NewSongController()
	r.POST("/song/add", songCtrl.AddSong)
	r.POST("/song/update", songCtrl.UpdateSong)
	r.GET("/song/delete", songCtrl.DeleteSong)
	r.GET("/song/detail", songCtrl.SongOfId)
	r.GET("/song/:id", songCtrl.SongOfId)
	r.GET("/song/singer/detail", songCtrl.SongOfSingerId)
	r.GET("/song/album/detail", songCtrl.SongsOfAlbumId) // 新增：查专辑歌曲列表
	r.GET("/song/name/detail", songCtrl.SongOfName)
	r.GET("/song", songCtrl.AllSong)

	// SongList routes
	songListCtrl := controller.NewSongListController()
	r.POST("/songList/add", songListCtrl.AddSongList)
	r.POST("/songList/update", songListCtrl.UpdateSongList)
	r.GET("/songList/delete", songListCtrl.DeleteSongList)
	r.GET("/songList/detail", songListCtrl.SongListOfId)
	r.GET("/songList/name/detail", songListCtrl.SongListOfTitle)
	r.GET("/songList", songListCtrl.AllSongList)

	// Collect routes
	collectCtrl := controller.NewCollectController()
	r.POST("/collect/add", collectCtrl.AddCollect)
	r.GET("/collect/delete", collectCtrl.DeleteCollect)
	r.GET("/collect/detail", collectCtrl.CollectOfUserId)

	// Comment routes
	commentCtrl := controller.NewCommentController()
	r.POST("/comment/add", commentCtrl.AddComment)
	r.GET("/comment/delete", commentCtrl.DeleteComment)
	r.GET("/comment/song/detail", commentCtrl.CommentOfSongId)
	r.GET("/comment/songList/detail", commentCtrl.CommentOfSongListId)

	// Banner routes
	bannerCtrl := controller.NewBannerController()
	r.GET("/banner", bannerCtrl.AllBanner)
	r.GET("/banner/getAllBanner", bannerCtrl.AllBanner) // 兼容旧路径

	// Website routes
	websiteCtrl := controller.NewWebsiteController()
	r.GET("/website", websiteCtrl.AllWebsite)

	// SingerRank routes - 排行榜歌手
	singerRankCtrl := controller.NewSingerRankController()
	r.GET("/singer-rank", singerRankCtrl.AllSinger)
	r.GET("/singer-rank/:id", singerRankCtrl.SingerOfId)
	r.GET("/singer-rank/name/detail", singerRankCtrl.SingerOfName)
	r.POST("/singer-rank/add", singerRankCtrl.AddSinger)
	r.POST("/singer-rank/update", singerRankCtrl.UpdateSinger)
	r.GET("/singer-rank/delete", singerRankCtrl.DeleteSinger)

	// AlbumRank routes - 排行榜专辑
	albumRankCtrl := controller.NewAlbumRankController()
	r.GET("/album-rank", albumRankCtrl.AllAlbum)
	r.GET("/album-rank/:id", albumRankCtrl.AlbumOfId)
	r.GET("/album-rank/singer/detail", albumRankCtrl.AlbumsOfSingerId)
	r.POST("/album-rank/add", albumRankCtrl.AddAlbum)
	r.POST("/album-rank/update", albumRankCtrl.UpdateAlbum)
	r.GET("/album-rank/delete", albumRankCtrl.DeleteAlbum)

	// Admin routes
	adminCtrl := controller.NewAdminController()
	r.POST("/admin/login", adminCtrl.Login)
	r.GET("/admin/login/status", adminCtrl.LoginStatus)

	// ListSong routes
	listSongCtrl := controller.NewListSongController()
	r.POST("/listSong/add", listSongCtrl.AddListSong)
	r.GET("/listSong/delete", listSongCtrl.DeleteListSong)
	r.GET("/listSong/detail", listSongCtrl.ListSongOfSongListId)

	// UserSupport routes
	userSupportCtrl := controller.NewUserSupportController()
	r.POST("/userSupport/add", userSupportCtrl.AddUserSupport)
	r.GET("/userSupport/delete", userSupportCtrl.DeleteUserSupport)

	// Import routes - 歌曲导入相关
	importCtrl := controller.NewImportController()
	r.POST("/import/format-name", importCtrl.FormatName)        // 名字格式化
	r.POST("/import/move-file", importCtrl.MoveFile)             // 移动文件到HDD
	r.POST("/import/songs", importCtrl.ImportSongs)              // 规整进数据库
	r.POST("/import/singer/albums", importCtrl.ImportSingerAlbums) // 一键导入-单歌手-所有专辑

	// Rank routes - 排行榜导入
	rankCtrl := controller.NewRankController()
	r.POST("/rank/check", rankCtrl.CheckRank)    // 校验排行榜数据（不入库）
	r.POST("/rank/import", rankCtrl.ImportRank)   // 导入排行榜数据
	r.GET("/rank/list", rankCtrl.GetRankList)    // 获取榜单列表
	r.GET("/rank/detail", rankCtrl.GetRankDetail) // 获取榜单详情

	// S3 routes - AWS S3 文件夹管理
	s3Ctrl := controller.NewS3Controller()
	r.GET("/s3/folders", s3Ctrl.ListFolders)               // 列出文件夹列表
	r.POST("/s3/folder", s3Ctrl.CreateFolder)              // 创建文件夹
	// r.DELETE("/s3/folder", s3Ctrl.DeleteFolder)         // 删除文件夹（已禁用，防止误操作）
	r.GET("/s3/folder/info", s3Ctrl.GetFolderInfo)         // 获取文件夹详情
	r.GET("/s3/objects", s3Ctrl.ListObjects)               // 列出文件夹内文件
	r.POST("/s3/upload", s3Ctrl.UploadFile)                // 上传文件
	r.GET("/s3/object/info", s3Ctrl.GetObjectInfo)         // 获取文件详情
	r.GET("/s3/object/download", s3Ctrl.DownloadObject)    // 获取文件下载链接
	r.POST("/s3/object/copy", s3Ctrl.CopyObject)           // 复制文件
	// r.DELETE("/s3/object", s3Ctrl.DeleteObject)          // 删除文件（已禁用，防止误操作）
	// 批量操作
	r.POST("/s3/objects/batch/info", s3Ctrl.BatchGetObjectsInfo)  // 批量获取文件信息
	r.POST("/s3/objects/batch/copy", s3Ctrl.BatchCopyObjects)     // 批量复制文件
	// r.DELETE("/s3/objects/batch", s3Ctrl.BatchDeleteObjects)  // 批量删除文件（已禁用，防止误操作）
}
