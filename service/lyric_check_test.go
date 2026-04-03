package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"study-music-server-go/config"
	"study-music-server-go/mapper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var checkDB *gorm.DB

func initCheckDB() error {
	// 获取项目根目录
	dirs := []string{".", "..", "../..", "../../.."}
	var projectRoot string
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(d, "config.yaml")); err == nil {
			abs, _ := filepath.Abs(d)
			projectRoot = abs
			break
		}
	}

	// 加载配置
	configPath := filepath.Join(projectRoot, "config.yaml")
	_, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// 连接数据库
	checkDB, err = gorm.Open(mysql.Open(config.AppConfig.Database.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 设置全局DB
	mapper.DB = checkDB
	return nil
}

// TestCheckLyricData 检查数据库中歌词数据情况
func TestCheckLyricData(t *testing.T) {
	// 初始化数据库连接
	if err := initCheckDB(); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 查询所有歌曲的歌词
	var songs []struct {
		ID    uint
		Name  string
		Lyric string
	}

	mapper.DB.Raw(`
		SELECT id, name, lyric
		FROM song
		WHERE lyric != ''
	`).Scan(&songs)

	fmt.Printf("共有 %d 首歌曲有歌词\n", len(songs))

	// 统计正常歌词和异常歌词
	var normalCount, abnormalCount int
	for _, song := range songs {
		// 去掉 BOM 字符后判断
		lyric := song.Lyric
		if len(lyric) > 0 && lyric[0] == 0xEF && len(lyric) > 2 && lyric[1] == 0xBB && lyric[2] == 0xBF {
			lyric = lyric[3:] // 去掉 UTF-8 BOM
		}
		isNormal := len(lyric) > 0 && lyric[0] == '['
		if isNormal {
			normalCount++
		} else {
			abnormalCount++
			// 打印异常歌词的前50个字符
			preview := song.Lyric
			if len(preview) > 50 {
				preview = preview[:50]
			}
			fmt.Printf("异常歌词: id=%d, name=%s, lyric_preview=%s\n", song.ID, song.Name, preview)
		}
	}

	fmt.Printf("\n统计: 正常歌词=%d, 异常歌词=%d\n", normalCount, abnormalCount)

	if abnormalCount > 0 {
		t.Logf("发现 %d 条异常歌词数据（不是lrc格式）", abnormalCount)
	}
}

// TestCleanAbnormalLyric 清理异常歌词数据（去掉 [awlrc:...] 部分）
func TestCleanAbnormalLyric(t *testing.T) {
	if err := initCheckDB(); err != nil {
		t.Fatalf("初始化数据库失败: %v", err)
	}

	// 查询所有包含 awlrc 的歌词
	var songs []struct {
		ID    uint
		Name  string
		Lyric string
	}
	mapper.DB.Raw(`SELECT id, name, lyric FROM song WHERE lyric LIKE '%awlrc%'`).Scan(&songs)

	fmt.Printf("包含 awlrc 的歌词数量: %d\n", len(songs))

	if len(songs) == 0 {
		t.Log("没有需要清理的歌词")
		return
	}

	// 清理每首歌的歌词
	for _, song := range songs {
		// 找到第一个 [awlrc 的位置，截取之前的内容
		idx := strings.Index(song.Lyric, "[awlrc")
		if idx > 0 {
			cleanLyric := strings.TrimSpace(song.Lyric[:idx])
			mapper.DB.Exec(`UPDATE song SET lyric = ? WHERE id = ?`, cleanLyric, song.ID)
			fmt.Printf("已清理: id=%d, name=%s\n", song.ID, song.Name)
		}
	}

	fmt.Printf("清理完成\n")
}