package service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"study-music-server-go/common"
	"study-music-server-go/mapper"
	"study-music-server-go/models"
	"study-music-server-go/utils"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type ImportService struct {
	singerMapper     *mapper.SingerMapper
	albumMapper      *mapper.AlbumMapper
	songMapper       *mapper.SongMapper
	songSingerMapper *mapper.SongSingerMapper
}

// readLrcFile 读取 lrc 文件，自动检测并转换编码（GBK -> UTF-8）
func readLrcFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	// 使用 Go 自带的 utf8.Valid 检测是否为有效的 UTF-8
	if utf8.Valid(content) {
		return string(content), nil
	}

	// UTF-8 无效，尝试 GBK 解码
	decoder := simplifiedchinese.GBK.NewDecoder()
	result, _, err := transform.Bytes(decoder, content)
	if err != nil {
		// 如果转换失败，返回原始内容
		return string(content), nil
	}

	return string(result), nil
}

func NewImportService() *ImportService {
	return &ImportService{
		singerMapper:     mapper.NewSingerMapper(),
		albumMapper:      mapper.NewAlbumMapper(),
		songMapper:       mapper.NewSongMapper(),
		songSingerMapper: mapper.NewSongSingerMapper(),
	}
}

// FormatName 名字格式化
// 指定"歌手-专辑"路径，遍历 mp3/wav/lrc 文件，重新格式化为"多作者-歌名.文件类型"
func (s *ImportService) FormatName(path string) *common.Response {
	files, err := utils.GetMusicFiles(path)
	if err != nil {
		return common.Error(fmt.Sprintf("读取目录失败: %v", err))
	}

	if len(files) == 0 {
		return common.Error("目录中没有音乐文件")
	}

	var results []map[string]string
	var failed []string

	for i := range files {
		file := &files[i]

		// 解析文件名获取歌手和歌曲名
		singer, songName := utils.ParseMusicFileName(file.OriginalName)

		// 格式化：多作者-歌名.扩展名
		file.NewName = utils.FormatMusicFileName(singer, songName, file.Ext)

		// 如果新名字和原名字不同，重命名文件
		if file.NewName != file.OriginalName {
			dir := filepath.Dir(file.Path)
			newPath := filepath.Join(dir, file.NewName)
			err = utils.MoveFile(file.Path, newPath)
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s -> %s: %v", file.OriginalName, file.NewName, err))
				continue
			}
			file.Path = newPath
		}

		results = append(results, map[string]string{
			"original": file.OriginalName,
			"new":      file.NewName,
		})
	}

	result := map[string]interface{}{
		"total":   len(files),
		"success": results,
		"failed":  failed,
	}

	return common.SuccessWithData("格式化完成", result)
}

// normalizeSMBPath 标准化SMB路径，处理 \\ 和 // 开头的情况
// Windows SMB: \\100.86.118.11\hdd -> \\100.86.118.11\hdd
// Linux SMB: //100.86.118.11/hdd -> //100.86.118.11/hdd
func normalizeSMBPath(path string) string {
	// 处理开头的 \\ 或 //
	if len(path) >= 2 {
		if path[0] == '\\' && path[1] == '\\' {
			// Windows SMB: \\server\share -> 保留
			return path
		}
		if path[0] == '/' && path[1] == '/' {
			// Linux SMB: //server/share -> 转换为 \\\\?\\UNC\\server\\share (Windows) 或保持原样
			// Go的filepath在处理//开头的路径时需要特殊处理
			// 统一转换为反斜杠格式以便filepath正确处理
			return strings.Replace(path, "/", "\\", -1)
		}
	}
	return path
}

// splitPathParts 分割路径，获取歌手名和专辑名
// 支持普通路径和SMB路径
// 注意：路径最后是专辑名，前一位是歌手名
func splitPathParts(path string) (singerName, albumName string) {
	// 先标准化SMB路径
	path = normalizeSMBPath(path)

	// 去除末尾的斜杠
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, "\\")

	// 使用反斜杠和正斜杠分割
	parts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '\\' || r == '/'
	})

	// 跳过开头可能出现的空元素（如 \\server\share 分割后第一个是空）
	// 找到第一个非空元素的索引
	startIdx := 0
	for i, p := range parts {
		if p != "" {
			startIdx = i
			break
		}
	}

	// 检测UNC路径（\\server\share），跳过服务器名和共享名
	// UNC路径格式：\\server\share\path1\path2
	if len(parts) > startIdx+2 {
		// 检查是否是UNC路径（开头是 \\ 或 //）
		if (len(path) >= 2 && path[0] == '\\' && path[1] == '\\') ||
			(len(path) >= 2 && path[0] == '/' && path[1] == '/') {
			// UNC路径：跳过 server 和 share，取后面的路径
			startIdx = startIdx + 2
		}
	}

	log.Printf("splitPathParts: path=%s, parts=%v, startIdx=%d", path, parts, startIdx)

	if len(parts)-startIdx >= 2 {
		albumName = parts[len(parts)-1]
		singerName = parts[len(parts)-2]
	}

	log.Printf("splitPathParts result: singerName=%s, albumName=%s", singerName, albumName)
	return singerName, albumName
}

// buildTargetPath 构建目标路径
func buildTargetPath(rootDir, singerName, albumName string) string {
	// 如果是SMB路径（以\\或//开头），需要特殊处理
	rootDir = normalizeSMBPath(rootDir)

	// 使用 filepath.Join 来正确拼接路径
	return filepath.Join(rootDir, singerName, albumName)
}

// MoveFile 移动文件到HDD
// from: 源目录路径，to: 目标根目录（会自动创建 歌手名/专辑名/ 子目录）
// 支持普通路径和SMB路径（Windows: \\server\share, Linux: //server/share）
// 只有全部文件都移动成功才返回成功，否则返回失败
func (s *ImportService) MoveFile(from, to string) *common.Response {
	// 检查源目录是否存在
	if !utils.FileExists(from) {
		return common.Error("源目录不存在: " + from)
	}

	// 获取源目录下的所有音乐文件
	files, err := utils.GetMusicFiles(from)
	if err != nil {
		return common.Error(fmt.Sprintf("读取目录失败: %v", err))
	}

	if len(files) == 0 {
		return common.Error("目录中没有音乐文件")
	}

	// 解析源路径，获取歌手名和专辑名（用于构建 NAS URL）
	singerName, albumName := splitPathParts(from)

	if singerName == "" || albumName == "" {
		return common.Error("无法从路径中解析歌手名和专辑名，请确保路径格式为：根目录/歌手名/专辑名/")
	}

	// 目标路径直接使用（用户填的就是完整路径：\\server\share\歌手\专辑）
	// 统一路径分隔符
	targetDir := strings.Replace(to, "/", "\\", -1)

	// 创建目标目录
	log.Printf("创建目标目录: %s", targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return common.Error(fmt.Sprintf("创建目标目录失败: %v", err))
	}
	log.Printf("目标目录创建成功")

	// 记录源文件信息（用于最后验证）
	var srcFileSizes map[string]int64
	srcFileSizes = make(map[string]int64)
	for _, file := range files {
		if info, err := os.Stat(file.Path); err == nil {
			srcFileSizes[file.OriginalName] = info.Size()
		}
	}

	var failed []string

	// 移动每个文件
	for i := range files {
		file := &files[i]
		targetPath := filepath.Join(targetDir, file.OriginalName)

		err := utils.MoveFile(file.Path, targetPath)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s -> %s: %v", file.OriginalName, targetPath, err))
			continue
		}

		// 验证移动是否成功
		if !utils.FileExists(targetPath) {
			failed = append(failed, fmt.Sprintf("%s: 移动后文件不存在", file.OriginalName))
			continue
		}
	}

	// 如果有失败，直接返回失败
	if len(failed) > 0 {
		return common.Error(fmt.Sprintf("移动失败，共 %d 个文件", len(failed)))
	}

	// 全部移动完成后，验证目标目录的文件数量和大小
	targetFiles, err := utils.GetMusicFiles(targetDir)
	if err != nil {
		return common.Error(fmt.Sprintf("验证失败：无法读取目标目录: %v", err))
	}

	// 检查文件数量
	if len(targetFiles) != len(files) {
		return common.Error(fmt.Sprintf("验证失败：源文件 %d 个，目标文件 %d 个", len(files), len(targetFiles)))
	}

	// 检查每个文件的大小
	for _, tf := range targetFiles {
		srcSize, ok := srcFileSizes[tf.OriginalName]
		if !ok {
			return common.Error(fmt.Sprintf("验证失败：目标文件 %s 在源文件中找不到", tf.OriginalName))
		}
		if info, err := os.Stat(tf.Path); err != nil {
			return common.Error(fmt.Sprintf("验证失败：无法获取目标文件 %s 信息: %v", tf.OriginalName, err))
		} else if info.Size() != srcSize {
			return common.Error(fmt.Sprintf("验证失败：文件 %s 大小不匹配，源: %d bytes, 目标: %d bytes", tf.OriginalName, srcSize, info.Size()))
		}
	}

	result := map[string]interface{}{
		"total":      len(files),
		"singer":     singerName,
		"album":      albumName,
		"target_dir": targetDir,
	}

	// 移动成功后，自动入库
	importResult := s.ImportSongs(targetDir)
	result["import_result"] = importResult.Data

	return common.SuccessWithData("移动成功", result)
}

// ImportSingerAlbums 一键导入-单歌手-所有专辑
// singerPath: 歌手目录路径，如 C:\test\周杰伦（目录下有多个专辑子目录）
// targetPath: 目标路径，如 \\100.86.118.11\hdd\周杰伦
func (s *ImportService) ImportSingerAlbums(singerPath, targetPath string) *common.Response {
	// 标准化路径
	singerPath = strings.Replace(singerPath, "/", "\\", -1)
	targetPath = strings.Replace(targetPath, "/", "\\", -1)

	// 从源路径中提取歌手名（源路径最后一级就是歌手名）
	singerName := filepath.Base(singerPath)
	if singerName == "." || singerName == "" {
		return common.Error("无法从路径中提取歌手名")
	}
	log.Printf("ImportSingerAlbums: singerName=%s, targetPath=%s", singerName, targetPath)

	// 确保SMB连接
	utils.EnsureSMBConnection(targetPath)

	// 读取源目录下的所有子目录（每个子目录是一个专辑）
	entries, err := os.ReadDir(singerPath)
	if err != nil {
		return common.Error(fmt.Sprintf("读取歌手目录失败: %v", err))
	}

	var results []map[string]interface{}
	var failed []map[string]interface{}
	var totalSongs int

	for _, entry := range entries {
		if !entry.IsDir() {
			continue // 跳过文件，只处理目录
		}

		albumName := entry.Name()
		albumPath := filepath.Join(singerPath, albumName)

		// 检查专辑目录下是否有音乐文件
		files, err := utils.GetMusicFiles(albumPath)
		if err != nil {
			failed = append(failed, map[string]interface{}{
				"album": albumName,
				"error": fmt.Sprintf("读取目录失败: %v", err),
			})
			continue
		}

		if len(files) == 0 {
			continue // 跳过空目录
		}

		// 第一步：格式化文件名
		log.Printf("格式化专辑文件名: %s", albumPath)
		formatResult := s.FormatName(albumPath)
		if !formatResult.Success {
			log.Printf("格式化失败: %s", formatResult.Message)
		}

		// 构建目标完整路径：目标路径 + 专辑名
		// 如: \\100.86.118.11\hdd\林俊杰\JJ陆
		albumTargetPath := filepath.Join(targetPath, albumName)

		// 直接移动到目标路径，传入歌手名（不再从源路径解析）
		moveResult := s.moveFileDirect(albumPath, albumTargetPath, singerName)

		if moveResult.Success {
			totalSongs += len(files)
			results = append(results, map[string]interface{}{
				"album":   albumName,
				"success": true,
				"songs":   len(files),
				"message": "成功",
			})
		} else {
			failed = append(failed, map[string]interface{}{
				"album":   albumName,
				"success": false,
				"error":   moveResult.Message,
			})
		}
	}

	response := map[string]interface{}{
		"total_albums":  len(results) + len(failed),
		"success":       len(results),
		"failed":        len(failed),
		"total_songs":   totalSongs,
		"results":       results,
		"failed_detail": failed,
	}

	// 如果没有导入任何歌曲，返回错误
	if totalSongs == 0 {
		return common.Error("源目录下没有音乐文件")
	}

	if len(failed) > 0 && len(results) == 0 {
		return common.Error("所有专辑导入失败")
	}

	return common.SuccessWithData("导入完成", response)
}

// moveFileDirect 直接移动文件到指定目标路径（不经过路径构建）
// from: 源目录路径，to: 目标完整路径（包含歌手名和专辑名）
func (s *ImportService) moveFileDirect(from, to, singerName string) *common.Response {
	log.Printf("moveFileDirect: from=%s, to=%s, singerName=%s", from, to, singerName)

	// 检查源目录是否存在
	if !utils.FileExists(from) {
		return common.Error("源目录不存在: " + from)
	}

	// 获取源目录下的所有音乐文件
	files, err := utils.GetMusicFiles(from)
	if err != nil {
		return common.Error(fmt.Sprintf("读取目录失败: %v", err))
	}

	if len(files) == 0 {
		return common.Error("目录中没有音乐文件")
	}

	// 解析源路径，获取歌手名和专辑名（用于NAS URL）
	// 解析专辑名（从源路径最后一级获取）
	albumName := filepath.Base(from)

	if albumName == "" {
		return common.Error("无法从路径中解析专辑名")
	}

	// 目标路径直接使用
	targetDir := strings.Replace(to, "/", "\\", -1)

	// 创建目标目录
	log.Printf("创建目标目录: %s", targetDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return common.Error(fmt.Sprintf("创建目标目录失败: %v", err))
	}
	log.Printf("目标目录创建成功")

	// 记录源文件信息
	var srcFileSizes map[string]int64
	srcFileSizes = make(map[string]int64)
	for _, file := range files {
		if info, err := os.Stat(file.Path); err == nil {
			srcFileSizes[file.OriginalName] = info.Size()
		}
	}

	var failed []string

	// 移动每个文件
	for i := range files {
		file := &files[i]
		targetPath := filepath.Join(targetDir, file.OriginalName)

		err := utils.MoveFile(file.Path, targetPath)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s -> %s: %v", file.OriginalName, targetPath, err))
			continue
		}

		// 验证移动是否成功
		if !utils.FileExists(targetPath) {
			failed = append(failed, fmt.Sprintf("%s: 移动后文件不存在", file.OriginalName))
			continue
		}
	}

	// 如果有失败，直接返回失败
	if len(failed) > 0 {
		return common.Error(fmt.Sprintf("移动失败，共 %d 个文件", len(failed)))
	}

	// 全部移动完成后，验证目标目录
	targetFiles, err := utils.GetMusicFiles(targetDir)
	if err != nil {
		return common.Error(fmt.Sprintf("验证失败：无法读取目标目录: %v", err))
	}

	if len(targetFiles) != len(files) {
		return common.Error(fmt.Sprintf("验证失败：源文件 %d 个，目标文件 %d 个", len(files), len(targetFiles)))
	}

	result := map[string]interface{}{
		"total":      len(files),
		"singer":     singerName,
		"album":      albumName,
		"target_dir": targetDir,
	}

	// 自动入库
	importResult := s.ImportSongs(targetDir)
	result["import_result"] = importResult.Data

	return common.SuccessWithData("移动成功", result)
}

// ImportSongs 规整进数据库
// 指定路径，遍历 mp3/wav/lrc 文件，插入/查询歌手、插入歌曲、关联关系
func (s *ImportService) ImportSongs(path string) *common.Response {
	files, err := utils.GetMusicFiles(path)
	if err != nil {
		return common.Error(fmt.Sprintf("读取目录失败: %v", err))
	}

	if len(files) == 0 {
		return common.Error("目录中没有音乐文件")
	}

	// 从路径中提取歌手名和专辑名
	// 路径格式：/path/歌手/专辑/ 或 /path/歌手-专辑/ 或 C:\path\歌手\专辑\
	// 使用 strings.FieldsFunc 同时支持 / 和 \ 分隔符
	path = strings.TrimSuffix(path, "/")
	path = strings.TrimSuffix(path, "\\")
	pathParts := strings.FieldsFunc(path, func(r rune) bool {
		return r == '/' || r == '\\'
	})
	var singerName, albumName string
	if len(pathParts) >= 2 {
		albumName = pathParts[len(pathParts)-1]
		singerName = pathParts[len(pathParts)-2]
	}

	// 处理歌手名中的分隔符（如"歌手-专辑"格式）
	if strings.Contains(singerName, "-") {
		parts := strings.Split(singerName, "-")
		singerName = parts[0]
		if len(parts) > 1 && albumName == "" {
			albumName = parts[1]
		}
	}

	var importedSongs []map[string]interface{}
	var failed []string

	for i := range files {
		file := &files[i]

		// 解析文件名获取歌手和歌曲名
		fileSinger, songName := utils.ParseMusicFileName(file.OriginalName)

		// 优先使用文件名中的歌手名，如果没有则使用目录名
		if fileSinger != "" {
			singerName = fileSinger
		}

		// 跳过歌词文件，不入库
		if file.Ext == ".lrc" {
			continue
		}

		// 处理多个歌手（按"、"分割，如"周杰伦、杨瑞代"）
		singerNames := strings.Split(singerName, "、")
		var singerIds []uint

		// 插入或查询歌手
		for _, name := range singerNames {
			name = strings.TrimSpace(name)
			if name == "" {
				continue
			}

			singer, err := s.singerMapper.FindByName(name)
			if err != nil || len(singer) == 0 {
				// 歌手不存在，插入新歌手
				newSinger := &models.Singer{Name: name}
				err = s.singerMapper.Add(newSinger)
				if err != nil {
					failed = append(failed, fmt.Sprintf("插入歌手失败: %s - %v", name, err))
					continue
				}
				singerIds = append(singerIds, newSinger.ID)
			} else {
				singerIds = append(singerIds, singer[0].ID)
			}
		}

		if len(singerIds) == 0 {
			failed = append(failed, fmt.Sprintf("无法获取歌手ID: %s", file.OriginalName))
			continue
		}

		// 处理专辑
		var albumId *uint
		if albumName != "" && len(singerIds) > 0 {
			// 使用第一个歌手的专辑
			album, err := s.albumMapper.FindByNameAndSingerId(albumName, singerIds[0])
			if err != nil || album == nil {
				// 专辑不存在，插入新专辑
				newAlbum := &models.Album{
					Name:     albumName,
					SingerId: singerIds[0],
				}
				err = s.albumMapper.Add(newAlbum)
				if err != nil {
					// 专辑插入失败，继续创建歌曲（album_id 为空）
					fmt.Printf("插入专辑失败: %v\n", err)
				} else {
					albumId = &newAlbum.ID
				}
			} else {
				albumId = &album.ID
			}
		}

		// 构建NAS路径：歌手名/专辑名/原始文件名
		firstSinger := singerNames[0]
		if firstSinger == "" {
			firstSinger = "未知歌手"
		}
		nasUrl := fmt.Sprintf("%s/%s/%s", firstSinger, albumName, file.OriginalName)

		// 只在多人时存储歌手名，单人则为空
		var fullNameSinger string
		if len(singerNames) > 1 {
			fullNameSinger = singerName
		}

		// 插入歌曲
		song := &models.Song{
			Name:           songName + file.Ext, // 歌曲名包含扩展名
			AlbumId:        albumId,
			NasUrlPath:     nasUrl,
			FullNameSinger: fullNameSinger, // 多歌手时存储，单人则为空
		}
		err = s.songMapper.Add(song)
		if err != nil {
			failed = append(failed, fmt.Sprintf("插入歌曲失败: %s - %v", songName, err))
			continue
		}

		// 插入歌曲歌手关联
		for _, singerId := range singerIds {
			songSinger := &models.SongSinger{
				SongId:   song.ID,
				SingerId: singerId,
			}
			err = s.songSingerMapper.Add(songSinger)
			if err != nil {
				fmt.Printf("插入歌曲歌手关联失败: %v\n", err)
			}
		}

		// 读取同名lrc文件作为歌词
		lrcPath := strings.TrimSuffix(file.Path, file.Ext) + ".lrc"
		fmt.Printf("检查歌词文件: %s\n", lrcPath)
		if lrcContent, err := readLrcFile(lrcPath); err == nil {
			fmt.Printf("找到歌词文件: %s, 大小: %d bytes\n", lrcPath, len(lrcContent))
			// 找到歌曲，更新歌词
			foundSong, err := s.songMapper.FindByName(song.Name)
			if err == nil && len(foundSong) > 0 {
				foundSong[0].Lyric = string(lrcContent)
				s.songMapper.Update(&foundSong[0])
				fmt.Printf("歌词已更新到歌曲: %s\n", song.Name)
			} else {
				fmt.Printf("未找到歌曲: %s, err: %v\n", song.Name, err)
			}
		} else {
			fmt.Printf("未找到歌词文件: %s, err: %v\n", lrcPath, err)
		}

		importedSongs = append(importedSongs, map[string]interface{}{
			"name":             song.Name,
			"singer":           singerName,
			"album":            albumName,
			"nas_url":          nasUrl,
			"full_name_singer": fullNameSinger,
			"song_id":          song.ID,
		})
	}

	result := map[string]interface{}{
		"total":    len(files),
		"imported": importedSongs,
		"failed":   failed,
	}

	return common.SuccessWithData("导入完成", result)
}
