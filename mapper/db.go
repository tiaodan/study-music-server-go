package mapper

import (
	"study-music-server-go/config"
	"study-music-server-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	cfg := config.AppConfig

	DB, err = gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logger.Info,
			},
		),
	})
	if err != nil {
		return err
	}

	// 自动迁移，创建或更新表结构
	err = DB.AutoMigrate(
	    &models.Consumer{},
	    &models.Singer{},
	    &models.Album{},
	    &models.Song{},
	    &models.SongList{},
	    &models.Collect{},
	    &models.Comment{},
	    &models.RankList{},
	    &models.Banner{},
	    &models.Admin{},
	    &models.ListSong{},
	    &models.UserSupport{},
	    &models.Device{},
	    &models.Website{},
	    &models.SongRank{},
	    &models.Rank{},
	)

	// 自动修复字符集问题
	fixCharset()

	// 初始化 website 预设数据
	initWebsiteData()

	return nil
}

// fixCharset 修复数据库表的字符集问题
func fixCharset() {
	// 先修复数据库默认字符集
	DB.Exec("ALTER DATABASE CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	// 修复 song 表的 lyric 字段字符集
	DB.Exec("ALTER TABLE song CONVERT TO CHARACTER SET utf8mb4")
	DB.Exec("ALTER TABLE song MODIFY lyric TEXT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	// 修复其他可能的中文字段
	DB.Exec("ALTER TABLE song MODIFY name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	DB.Exec("ALTER TABLE singer MODIFY name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	DB.Exec("ALTER TABLE album MODIFY name VARCHAR(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")

	log.Println("数据库字符集已修复为 utf8mb4")
}

// initWebsiteData 初始化 website 预设数据
func initWebsiteData() {
	websites := []models.Website{
		{Name: "QQ音乐", Type: "music"},
		{Name: "酷狗音乐", Type: "music"},
		{Name: "酷我音乐", Type: "music"},
		{Name: "网易云音乐", Type: "music"},
		{Name: "咪咕音乐", Type: "music"},
	}
	for _, w := range websites {
		DB.FirstOrCreate(&w, models.Website{Name: w.Name, Type: w.Type})
	}
	log.Println("Website预设数据已初始化")
}
