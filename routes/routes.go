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

	// Song routes
	songCtrl := controller.NewSongController()
	r.POST("/song/add", songCtrl.AddSong)
	r.POST("/song/update", songCtrl.UpdateSong)
	r.GET("/song/delete", songCtrl.DeleteSong)
	r.GET("/song/detail", songCtrl.SongOfId)
	r.GET("/song/:id", songCtrl.SongOfId)
	r.GET("/song/singer/detail", songCtrl.SongOfSingerId)
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

	// RankList routes
	rankListCtrl := controller.NewRankListController()
	r.POST("/rankList/add", rankListCtrl.AddRankList)
	r.GET("/rankList/detail", rankListCtrl.RankListOfSongListId)

	// Banner routes
	bannerCtrl := controller.NewBannerController()
	r.GET("/banner", bannerCtrl.AllBanner)

	// Admin routes
	adminCtrl := controller.NewAdminController()
	r.POST("/admin/login", adminCtrl.Login)

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
	r.POST("/import/format-name", importCtrl.FormatName)   // 名字格式化
	r.POST("/import/move-file", importCtrl.MoveFile)       // 移动文件到HDD
	r.POST("/import/songs", importCtrl.ImportSongs)        // 规整进数据库
}
