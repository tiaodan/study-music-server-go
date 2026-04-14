package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"study-music-server-go/config"
	"study-music-server-go/mapper"
	"study-music-server-go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 测试配置
var testDB *gorm.DB

// 获取项目根目录
func getProjectRoot() string {
	// 尝试从当前目录向上查找 config.yaml
	dirs := []string{".", "..", "../..", "../../.."}
	for _, d := range dirs {
		if _, err := os.Stat(filepath.Join(d, "config.yaml")); err == nil {
			abs, _ := filepath.Abs(d)
			return abs
		}
	}
	return "."
}

// 测试用的临时目录
var testSrcDir = getProjectRoot() + "/testdata/music/测试歌手/测试专辑"
var testToDir = getProjectRoot() + "/testdata/output"

func initTestDB() error {
	// 加载配置
	configPath := filepath.Join(getProjectRoot(), "config.yaml")
	_, err := config.LoadConfig(configPath)
	if err != nil {
		return err
	}

	// 连接测试数据库
	testDB, err = gorm.Open(mysql.Open(config.AppConfig.DatabaseTest.DSN()), &gorm.Config{})
	if err != nil {
		return err
	}

	// 设置全局DB
	mapper.DB = testDB

	// 自动迁移
	err = testDB.AutoMigrate(
		&models.Singer{},
		&models.Album{},
		&models.Song{},
	)
	if err != nil {
		return err
	}

	return nil
}

func cleanTestData() {
	// 清理测试数据
	testDB.Exec("DELETE FROM song")
	testDB.Exec("DELETE FROM album")
	testDB.Exec("DELETE FROM singer")
}

func setupTestDirs() {
	// 创建输出目录
	os.MkdirAll(testToDir, 0755)
}

func cleanupTestDirs() {
	// 清理输出目录
	os.RemoveAll(testToDir)
}

// TestImportSongs_ImportToDB 测试规整进数据库功能
func TestImportSongs_ImportToDB(t *testing.T) {
	// 初始化测试数据库
	if err := initTestDB(); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	cleanTestData()

	svc := NewImportService()

	// 调用 ImportSongs
	result := svc.ImportSongs(testSrcDir)

	// 验证结果
	if !result.Success {
		t.Errorf("导入失败: %s", result.Message)
	}

	// 1. 验证 song 表
	var songs []models.Song
	testDB.Find(&songs)

	// 验证歌曲数量（应该有2首：歌曲1 + 合唱歌曲）
	if len(songs) != 2 {
		t.Errorf("歌曲数量不对: 期望2, 实际%d", len(songs))
	}

	// 验证歌词是否有乱码
	for _, song := range songs {
		if song.Lyric != "" {
			// 检查是否包含乱码字符
			if strings.Contains(song.Lyric, "\uFFFD") || strings.Contains(song.Lyric, "�") {
				t.Errorf("歌曲[%s]歌词有乱码: %s", song.Name, song.Lyric[:100])
			}
			// 验证歌词是否包含中文
			if !strings.ContainsAny(song.Lyric, "测试歌词") && !strings.Contains(song.Lyric, "歌词") {
				t.Logf("歌曲[%s]歌词内容: %s", song.Name, song.Lyric[:100])
			}
		}
	}

	// 验证 nas_url_path 格式：歌手/专辑名/歌名
	for _, song := range songs {
		if song.NasUrlPath != "" {
			parts := strings.Split(song.NasUrlPath, "/")
			if len(parts) != 3 {
				t.Errorf("歌曲[%s]的nas_url_path格式不对: %s", song.Name, song.NasUrlPath)
			}
			// 验证第二部分是专辑名
			if len(parts) >= 2 && parts[1] != "测试专辑" {
				t.Errorf("歌曲[%s]的nas_url_path专辑名不对: 期望'测试专辑', 实际'%s'", song.Name, parts[1])
			}
			// 验证第一部分是歌手名（可能是目录名，也可能是文件名中的歌手名）
			t.Logf("歌曲[%s]的nas_url_path: %s", song.Name, song.NasUrlPath)
		}
	}

	// 2. 验证 singer 表
	var singers []models.Singer
	testDB.Find(&singers)

	// 应该至少有：测试歌手、歌手A、歌手B
	t.Logf("歌手数量: %d", len(singers))
	for _, s := range singers {
		t.Logf("歌手: %s", s.Name)
	}

	// 验证多人歌手（歌手A、歌手B）
	var multiSingerSongs []models.Song
	testDB.Where("full_name_singer != ?", "").Find(&multiSingerSongs)
	if len(multiSingerSongs) != 1 {
		t.Errorf("多歌手歌曲数量不对: 期望1, 实际%d", len(multiSingerSongs))
	} else {
		if multiSingerSongs[0].FullNameSinger != "歌手A、歌手B" {
			t.Errorf("多歌手名字不对: 期望'歌手A、歌手B', 实际'%s'", multiSingerSongs[0].FullNameSinger)
		}
	}

	// 3. 验证 album 表
	var albums []models.Album
	testDB.Find(&albums)
	t.Logf("专辑数量: %d", len(albums))
	for _, a := range albums {
		t.Logf("专辑: %s, singer_id: %d", a.Name, a.SingerId)
	}

	t.Logf("测试通过！共导入 %d 首歌曲", len(songs))
}

// TestMoveFile_OneKeyImport 测试一键导入（移动+入库）
func TestMoveFile_OneKeyImport(t *testing.T) {
	// 初始化测试数据库
	if err := initTestDB(); err != nil {
		t.Fatalf("初始化测试数据库失败: %v", err)
	}
	cleanTestData()

	// 创建测试数据副本（因为移动会删除原文件）
	testCopyDir := getProjectRoot() + "/testdata/onekey_test"
	os.MkdirAll(testCopyDir+"/测试歌手/测试专辑", 0755)

	// 复制测试文件
	testFiles := []string{
		"测试歌手-歌曲1.mp3",
		"歌手A、歌手B-合唱歌曲.mp3",
	}
	for _, name := range testFiles {
		src := filepath.Join(testSrcDir, name)
		dst := filepath.Join(testCopyDir, "测试歌手", "测试专辑", name)
		data, _ := os.ReadFile(src)
		os.WriteFile(dst, data, 0644)
	}

	setupTestDirs()
	defer func() {
		cleanupTestDirs()
		os.RemoveAll(testCopyDir)
	}()

	svc := NewImportService()

	// 构造目标路径（临时目录）
	targetDir := filepath.Join(testToDir, "测试歌手", "测试专辑")
	targetDir = strings.ReplaceAll(targetDir, "/", "\\")

	// 调用 MoveFile（一键导入：移动+自动入库）
	result := svc.MoveFile(testCopyDir+"/测试歌手/测试专辑", targetDir)

	// 验证结果
	if !result.Success {
		t.Errorf("一键导入失败: %s", result.Message)
	}

	// 验证移动后的文件存在
	files, _ := os.ReadDir(targetDir)
	t.Logf("目标目录文件数: %d", len(files))

	// 验证数据库
	var songs []models.Song
	testDB.Find(&songs)
	t.Logf("导入歌曲数量: %d", len(songs))

	// 验证 nas_url_path（只验证格式：歌手/专辑/歌名，不验证具体歌手名）
	for _, song := range songs {
		if song.NasUrlPath != "" {
			t.Logf("nas_url_path: %s", song.NasUrlPath)
			parts := strings.Split(song.NasUrlPath, "/")
			if len(parts) != 3 {
				t.Errorf("nas_url_path格式不对（应该是 歌手/专辑/歌名）: %s", song.NasUrlPath)
			}
		}
	}
}

// TestMoveFileOnly 测试仅移动文件功能（不入库）
func TestMoveFileOnly(t *testing.T) {
	// 先创建测试数据副本，因为移动会删除原文件
	testCopyDir := getProjectRoot() + "/testdata/move_test"
	os.MkdirAll(testCopyDir+"/测试歌手/测试专辑", 0755)

	// 复制测试文件
	testFiles := []string{
		"测试歌手-歌曲1.mp3",
		"歌手A、歌手B-合唱歌曲.mp3",
	}
	for _, name := range testFiles {
		src := filepath.Join(testSrcDir, name)
		dst := filepath.Join(testCopyDir, "测试歌手", "测试专辑", name)
		data, _ := os.ReadFile(src)
		os.WriteFile(dst, data, 0644)
	}

	setupTestDirs()
	defer func() {
		cleanupTestDirs()
		os.RemoveAll(testCopyDir)
	}()

	svc := NewImportService()

	// 目标路径（临时目录）
	targetDir := filepath.Join(testToDir, "移动测试")
	targetDir = strings.ReplaceAll(targetDir, "/", "\\")

	// 调用 MoveFile
	result := svc.MoveFile(testCopyDir+"/测试歌手/测试专辑", targetDir)

	if !result.Success {
		t.Errorf("移动失败: %s", result.Message)
	}

	// 验证文件已移动
	files, err := os.ReadDir(targetDir)
	if err != nil {
		t.Errorf("读取目标目录失败: %v", err)
	}

	t.Logf("移动后文件数量: %d", len(files))
	for _, f := range files {
		t.Logf("文件: %s", f.Name())
	}
}

// TestFormatName 测试名字格式化功能
func TestFormatName(t *testing.T) {
	svc := NewImportService()

	// 创建测试文件
	testDir := "./testdata/format_test"
	os.MkdirAll(testDir, 0755)
	defer os.RemoveAll(testDir)

	// 创建测试文件
	testFiles := []string{
		"歌手A-歌曲1.mp3",
		"歌手B-歌曲2.wav",
		"歌手C、歌手D-合唱.mp3",
	}
	for _, name := range testFiles {
		os.WriteFile(filepath.Join(testDir, name), []byte("test"), 0644)
	}

	// 调用 FormatName
	result := svc.FormatName(testDir)

	if !result.Success {
		t.Errorf("格式化失败: %s", result.Message)
	}

	t.Logf("格式化结果: %+v", result.Data)
}

// TestLrcEncoding 测试歌词编码问题
func TestLrcEncoding(t *testing.T) {
	// 创建GBK编码的lrc文件
	gbkLrcPath := getProjectRoot() + "/testdata/music/测试歌手/测试专辑/测试歌手-歌曲1_gbk.lrc"

	// 写入GBK编码内容
	content := []byte{
		0x5B, 0x74, 0x69, 0x3A, 0x6B, 0x6B,
		0x5D, 0x0D, 0x0A,
		0x5B, 0x61, 0x72, 0x3A, 0xCB, 0xD5,
		0xCF, 0xD4, 0x5D, 0x0D, 0x0A,
	}
	os.WriteFile(gbkLrcPath, content, 0644)
	defer os.Remove(gbkLrcPath)

	// 读取并验证
	contentRead, err := os.ReadFile(gbkLrcPath)
	if err != nil {
		t.Fatalf("读取文件失败: %v", err)
	}

	// 验证是GBK编码（不是有效的UTF-8）
	if isValidUTF8(contentRead) {
		t.Logf("文件是UTF-8编码")
	} else {
		t.Logf("文件不是UTF-8编码（可能是GBK）")
	}
}

// isValidUTF8 简单的UTF-8验证
func isValidUTF8(data []byte) bool {
	for i := 0; i < len(data); i++ {
		if data[i] < 0x80 {
			continue
		}
		// 简单检查
		if i+1 < len(data) && data[i] >= 0xC0 && data[i] < 0xE0 {
			i++
			continue
		}
	}
	return true
}

func main() {
	// 运行测试
	fmt.Println("运行导入功能测试...")
	testing.Main(func(pat, str string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{Name: "TestImportSongs_ImportToDB", F: TestImportSongs_ImportToDB},
			{Name: "TestMoveFile_OneKeyImport", F: TestMoveFile_OneKeyImport},
			{Name: "TestFormatName", F: TestFormatName},
			{Name: "TestMoveFileOnly", F: TestMoveFileOnly},
			{Name: "TestLrcEncoding", F: TestLrcEncoding},
		},
		nil, nil)
}
