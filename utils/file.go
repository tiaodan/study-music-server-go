package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SaveFile(srcFile, folder string) string {
	// Simplified file save - returns filename
	// In production, implement proper file upload handling
	ext := ".jpg"
	filename := fmt.Sprintf("%d%s%s", time.Now().Unix(), randString(6), ext)

	// Create directory if not exists
	dir := filepath.Join(".", "img", folder)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	return filename
}

func randString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GetFileMD5(file *os.File) string {
	md5Hash := md5.New()
	_, err := io.Copy(md5Hash, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(md5Hash.Sum(nil))
}

func GetFileExt(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return filename[idx:]
}

// MusicFileInfo 音乐文件信息
type MusicFileInfo struct {
	OriginalName string // 原始文件名
	NewName      string // 格式化后的文件名
	Path         string // 文件完整路径
	Ext          string // 文件扩展名
	Singer       string // 歌手名（可能多个，用逗号分隔）
	SongName     string // 歌曲名
	Album        string // 专辑名
}

// GetMusicFiles 遍历目录获取音乐文件
func GetMusicFiles(dirPath string) ([]MusicFileInfo, error) {
	var files []MusicFileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(GetFileExt(info.Name()))
		if ext == ".mp3" || ext == ".wav" || ext == ".lrc" {
			files = append(files, MusicFileInfo{
				OriginalName: info.Name(),
				Path:         path,
				Ext:          ext,
			})
		}
		return nil
	})

	return files, err
}

// ParseMusicFileName 解析文件名，提取歌手和歌曲名
// 支持格式：多歌手-歌曲名.mp3 或 歌手-歌曲名.mp3
// 返回：歌手列表, 歌曲名
func ParseMusicFileName(filename string) (string, string) {
	// 去掉扩展名
	ext := GetFileExt(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// 按 "-" 分割
	parts := strings.Split(nameWithoutExt, "-")
	if len(parts) < 2 {
		// 没有分隔符，整个作为歌曲名
		return "", nameWithoutExt
	}

	// 最后一个部分是歌曲名
	songName := strings.TrimSpace(parts[len(parts)-1])
	singer := strings.TrimSpace(strings.Join(parts[:len(parts)-1], ","))

	return singer, songName
}

// FormatMusicFileName 格式化音乐文件名
// 格式：多作者-歌名.文件类型
func FormatMusicFileName(singer, songName, ext string) string {
	return fmt.Sprintf("%s-%s%s", singer, songName, ext)
}

// MoveFile 移动文件
// 支持本地文件和SMB/UNC路径
// 如果目标文件已存在且大小一致，则跳过（删除源文件）
// 如果目标文件已存在但大小不一致，则覆盖
func MoveFile(from, to string) error {
	// 创建目标目录
	toDir := filepath.Dir(to)
	if _, err := os.Stat(toDir); os.IsNotExist(err) {
		err = os.MkdirAll(toDir, 0755)
		if err != nil {
			return err
		}
	}

	// 检查源文件是否存在
	fromInfo, err := os.Stat(from)
	if err != nil {
		return fmt.Errorf("源文件不存在: %s", from)
	}

	// 检查目标文件是否已存在
	skip := false
	var toInfo os.FileInfo
	if toInfo, err = os.Stat(to); err == nil {
		// 目标文件已存在，比较大小
		if fromInfo.Size() == toInfo.Size() {
			// 大小一致，跳过（删除源文件）
			if err := os.Remove(from); err != nil {
				return fmt.Errorf("删除源文件失败: %v", err)
			}
			return nil
		}
		// 大小不一致，覆盖
		skip = true
	}

	// 尝试使用 os.Rename（同一文件系统内有效）
	err = os.Rename(from, to)
	if err == nil {
		return nil
	}

	// 如果目标文件已存在（且大小不一致），先删除目标文件再复制
	if skip {
		os.Remove(to)
	}

	// 如果 Rename 失败（可能是跨文件系统/SMB），使用复制+删除的方式
	// 复制文件
	fromFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}

	toFile, err := os.Create(to)
	if err != nil {
		fromFile.Close()
		return fmt.Errorf("创建目标文件失败: %v", err)
	}

	_, err = io.Copy(toFile, fromFile)
	if err != nil {
		fromFile.Close()
		toFile.Close()
		return fmt.Errorf("复制文件失败: %v", err)
	}

	// 同步写入
	err = toFile.Sync()
	if err != nil {
		fromFile.Close()
		toFile.Close()
		return fmt.Errorf("同步文件失败: %v", err)
	}

	// 验证目标文件大小（必须在关闭文件之前）
	fromInfo, err = fromFile.Stat()
	if err != nil {
		fromFile.Close()
		toFile.Close()
		return fmt.Errorf("获取源文件信息失败: %v", err)
	}
	toInfo, err = toFile.Stat()
	if err != nil {
		fromFile.Close()
		toFile.Close()
		return fmt.Errorf("获取目标文件信息失败: %v", err)
	}

	// 关闭文件
	fromFile.Close()
	toFile.Close()

	// 验证文件大小
	if fromInfo.Size() != toInfo.Size() {
		return fmt.Errorf("文件大小不匹配，源文件: %d bytes, 目标文件: %d bytes", fromInfo.Size(), toInfo.Size())
	}

	// 删除源文件
	err = os.Remove(from)
	if err != nil {
		return fmt.Errorf("删除源文件失败: %v", err)
	}

	return nil
}

// FileExists 检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
