package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SingerRank 歌手表
type SingerRank struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:100"`
}

func (SingerRank) TableName() string {
	return "singer_rank"
}

// SongRank 歌曲表
type SongRank struct {
	ID             uint   `gorm:"primaryKey"`
	SingerId       *uint  `gorm:"column:singer_id"`
	Name           string `gorm:"size:100"`
	FullNameSinger string `gorm:"size:255;column:full_name_singer"`
	Lyric          string `gorm:"type:text"`
}

func (SongRank) TableName() string {
	return "song_rank"
}

// FileInfo 文件信息
type FileInfo struct {
	Path         string
	OriginalName string
	Singer       string
	SongName     string
	Ext          string
}

// readLrcFile 读取 lrc 文件
func readLrcFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	var result string
	if utf8.Valid(content) {
		result = string(content)
	} else {
		decoder := simplifiedchinese.GBK.NewDecoder()
		decoded, _, err := transform.Bytes(decoder, content)
		if err != nil {
			result = string(content)
		} else {
			result = string(decoded)
		}
	}

	return cleanLrcContent(result), nil
}

// cleanLrcContent 清理歌词
func cleanLrcContent(lyric string) string {
	if len(lyric) > 2 && lyric[0] == 0xEF && lyric[1] == 0xBB && lyric[2] == 0xBF {
		lyric = lyric[3:]
	}

	nonStandardTags := []string{"[awlrc", "[krc", "[qlrc"}
	for _, tag := range nonStandardTags {
		idx := strings.Index(lyric, tag)
		if idx > 0 {
			lyric = strings.TrimSpace(lyric[:idx])
			break
		}
	}

	return lyric
}

// removeBrackets 去括号
func removeBrackets(name string) string {
	re := strings.NewReplacer("（", "(", "）", ")")
	name = re.Replace(name)

	result := ""
	inBracket := false
	for _, ch := range name {
		if ch == '(' {
			inBracket = true
			continue
		}
		if ch == ')' {
			inBracket = false
			continue
		}
		if !inBracket {
			result += string(ch)
		}
	}

	return strings.TrimSpace(result)
}

// getMusicFiles 获取音乐文件
func getMusicFiles(folder string) ([]FileInfo, error) {
	var files []FileInfo

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if ext != ".mp3" && ext != ".wav" {
			continue // 只索引 mp3/wav
		}

		// 解析 "歌手 - 歌名" 格式
		baseName := strings.TrimSuffix(name, ext)
		idx := strings.Index(baseName, " - ")
		var singer, songName string
		if idx > 0 {
			singer = strings.TrimSpace(baseName[:idx])
			songName = strings.TrimSpace(baseName[idx+3:])
		}

		files = append(files, FileInfo{
			Path:         filepath.Join(folder, name),
			OriginalName: name,
			Singer:       singer,
			SongName:     songName,
			Ext:          ext,
		})
	}

	return files, nil
}

// matchFile 匹配文件（模糊匹配）
func matchFile(files []FileInfo, singerName, songName string) *FileInfo {
	singerName = strings.TrimSpace(singerName)
	songName = strings.TrimSpace(songName)

	// 去括号版本
	songNameNoBracket := removeBrackets(songName)

	for _, f := range files {
		fSinger := strings.TrimSpace(f.Singer)
		fSongName := strings.TrimSpace(f.SongName)

		// 检查歌手是否匹配（支持多歌手）
		// 文件名可能是 "Justin Bieber、Ludacris"，数据库里是 "Justin Bieber"
		singerMatch := false
		if singerName != "" {
			// 完全匹配
			if fSinger == singerName {
				singerMatch = true
			}
			// 文件名包含数据库歌手名
			if strings.Contains(fSinger, singerName) {
				singerMatch = true
			}
			// 数据库歌手名包含文件名歌手（多歌手情况）
			if singerName != "" && fSinger != "" && strings.Contains(singerName, fSinger) {
				singerMatch = true
			}
			// 尝试第一个歌手（多歌手用顿号分隔）
			if strings.Contains(singerName, "、") {
				firstSinger := strings.Split(singerName, "、")[0]
				if strings.Contains(fSinger, firstSinger) {
					singerMatch = true
				}
			}
			if strings.Contains(fSinger, "、") {
				firstFileSinger := strings.Split(fSinger, "、")[0]
				if strings.Contains(singerName, firstFileSinger) {
					singerMatch = true
				}
			}
		}

		// 检查歌名是否匹配
		songMatch := false
		fSongNameNoBracket := removeBrackets(fSongName)

		// 完全匹配
		if fSongName == songName || fSongName == songNameNoBracket {
			songMatch = true
		}
		if fSongNameNoBracket == songName || fSongNameNoBracket == songNameNoBracket {
			songMatch = true
		}
		// 去空格后匹配
		if strings.ReplaceAll(fSongName, " ", "") == strings.ReplaceAll(songName, " ", "") {
			songMatch = true
		}

		if singerMatch && songMatch {
			return &f
		}
	}

	return nil
}

func main() {
	// 数据库连接
	dsn := "root:password@tcp(127.0.0.1:3306)/music?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	fmt.Println("数据库连接成功")

	// 文件夹路径
	folderPath := "C:/A_software/D_软件的保存路径/B. 电脑必备/04. 音乐/lx-download"

	// 读取文件夹
	files, err := getMusicFiles(folderPath)
	if err != nil {
		log.Fatalf("读取文件夹失败: %v", err)
	}
	fmt.Printf("文件夹下有 %d 个音乐文件\n", len(files))

	// 查询所有 song_rank（lyric 为空的）
	var songs []SongRank
	err = db.Where("lyric = '' OR lyric IS NULL").Find(&songs).Error
	if err != nil {
		log.Fatalf("查询歌曲失败: %v", err)
	}
	fmt.Printf("需要更新歌词的歌曲: %d 首\n", len(songs))

	// 查询所有歌手
	var singers []SingerRank
	db.Find(&singers)
	singerMap := make(map[uint]string)
	for _, s := range singers {
		singerMap[s.ID] = s.Name
	}

	var updatedCount int
	var failedCount int
	var failedList []string

	for _, song := range songs {
		// 获取歌手名
		singerName := ""
		if song.SingerId != nil {
			singerName = singerMap[*song.SingerId]
		}
		if singerName == "" && song.FullNameSinger != "" {
			singerName = song.FullNameSinger
		}

		// 匹配文件
		file := matchFile(files, singerName, song.Name)
		if file == nil {
			failedList = append(failedList, fmt.Sprintf("未找到: %s - %s", singerName, song.Name))
			failedCount++
			continue
		}

		// 找对应的 lrc 文件
		lrcPath := strings.TrimSuffix(file.Path, file.Ext) + ".lrc"
		lyric, err := readLrcFile(lrcPath)
		if err != nil {
			failedList = append(failedList, fmt.Sprintf("读取歌词失败: %s - %s (%v)", singerName, song.Name, err))
			failedCount++
			continue
		}

		// 更新数据库
		err = db.Model(&SongRank{}).Where("id = ?", song.ID).Update("lyric", lyric).Error
		if err != nil {
			failedList = append(failedList, fmt.Sprintf("更新失败: id=%d, err=%v", song.ID, err))
			failedCount++
			continue
		}

		updatedCount++
		log.Printf("更新成功: %s - %s (id=%d, 文件=%s, lyric=%d bytes)", singerName, song.Name, song.ID, file.OriginalName, len(lyric))
	}

	fmt.Printf("\n===== 完成 =====\n")
	fmt.Printf("成功更新: %d 首\n", updatedCount)
	fmt.Printf("失败: %d 首\n", failedCount)
	if len(failedList) > 0 {
		fmt.Println("失败列表:")
		for _, f := range failedList {
			fmt.Println("  ", f)
		}
	}
}
