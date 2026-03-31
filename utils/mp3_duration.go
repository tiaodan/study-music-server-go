package utils

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hajimehoshi/go-mp3"
)

// GetMP3Duration 读取 MP3 文件时长（秒）
func GetMP3Duration(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return 0, fmt.Errorf("解码失败: %v", err)
	}

	// 获取采样率和总采样数
	sampleRate := decoder.SampleRate()
	if sampleRate == 0 {
		return 0, fmt.Errorf("采样率为0")
	}

	// 计算时长：总采样数 / 采样率
	// decoder.Length() 返回的是PCM数据的字节数
	// 每个采样点占2字节（16位），双声道
	length := decoder.Length()
	if length == 0 {
		// 如果无法获取长度，尝试读取来计算
		buf := make([]byte, 4096)
		totalBytes := 0
		for {
			n, err := decoder.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return 0, fmt.Errorf("读取失败: %v", err)
			}
			totalBytes += n
		}
		length = int64(totalBytes)
	}

	// 时长 = 总字节数 / (采样率 * 2 * 2)
	// 采样率 * 2(字节) * 2(声道) = 每秒字节数
	duration := float64(length) / float64(sampleRate*4)

	return int(duration), nil
}

// IsAudioFile 判断是否为音频文件
func IsAudioFile(filename string) bool {
	ext := strings.ToLower(filename)
	return strings.HasSuffix(ext, ".mp3") || strings.HasSuffix(ext, ".wav") || strings.HasSuffix(ext, ".flac")
}