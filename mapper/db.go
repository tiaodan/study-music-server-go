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
	    &models.SongSinger{},
	    &models.SongList{},
	    &models.Collect{},
	    &models.Comment{},
	    &models.RankList{},
	    &models.Banner{},
	    &models.Admin{},
	    &models.ListSong{},
	    &models.UserSupport{},
	    &models.Device{},
	)

	return nil
}
