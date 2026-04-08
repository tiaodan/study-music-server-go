package service

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"study-music-server-go/common"
	"study-music-server-go/models"
	"study-music-server-go/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service struct {
	client    *s3.Client
	presigner *s3.PresignClient
	bucket    string
}

func NewS3Service() *S3Service {
	return &S3Service{
		client:    utils.S3Client,
		presigner: utils.S3Presigner,
		bucket:    utils.GetS3Bucket(),
	}
}

// ListFolders 列出文件夹列表
func (s *S3Service) ListFolders(prefix string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	// 规范化前缀
	prefix = normalizePath(prefix)

	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("列出文件夹失败: %v", err))
	}

	folders := []models.FolderInfo{}
	for _, prefix := range result.CommonPrefixes {
		folderPath := *prefix.Prefix
		folderName := getFolderName(folderPath)

		// 获取文件夹统计信息
		info := s.getFolderStats(folderPath)
		folders = append(folders, models.FolderInfo{
			Path:         folderPath,
			Name:         folderName,
			FileCount:    info.FileCount,
			TotalSize:    info.TotalSize,
			LastModified: info.LastModified,
		})
	}

	return common.SuccessWithData("获取成功", &models.ListFoldersResponse{
		Folders: folders,
		Total:   len(folders),
	})
}

// CreateFolder 创建文件夹
func (s *S3Service) CreateFolder(path string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	// 规范化路径（确保以 "/" 结尾）
	path = normalizePath(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// 检查是否已存在
	exists, err := s.folderExists(path)
	if err != nil {
		return common.Error(fmt.Sprintf("检查文件夹失败: %v", err))
	}
	if exists {
		return common.Error("文件夹已存在")
	}

	// 创建空对象（以 "/" 结尾的 key）
	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("创建文件夹失败: %v", err))
	}

	return common.Success("创建成功")
}

// DeleteFolder 删除文件夹
func (s *S3Service) DeleteFolder(path string, recursive bool) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	path = normalizePath(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// 列出文件夹下所有对象
	objects, err := s.listAllObjects(path)
	if err != nil {
		return common.Error(fmt.Sprintf("列出对象失败: %v", err))
	}

	// 如果不递归且文件夹非空
	if !recursive && len(objects) > 0 {
		return common.Error("文件夹非空，请使用递归删除")
	}

	// 批量删除（每次最多 1000 个）
	for len(objects) > 0 {
		batch := objects
		if len(objects) > 1000 {
			batch = objects[:1000]
			objects = objects[1000:]
		} else {
			objects = nil
		}

		deleteObjects := []types.ObjectIdentifier{}
		for _, obj := range batch {
			deleteObjects = append(deleteObjects, types.ObjectIdentifier{Key: obj.Key})
		}

		_, err = s.client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &types.Delete{
				Objects: deleteObjects,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return common.Error(fmt.Sprintf("删除失败: %v", err))
		}
	}

	return common.Success("删除成功")
}

// GetFolderInfo 获取文件夹详情
func (s *S3Service) GetFolderInfo(path string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	path = normalizePath(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	info := s.getFolderStats(path)
	folderName := getFolderName(path)

	return common.SuccessWithData("获取成功", &models.FolderInfo{
		Path:         path,
		Name:         folderName,
		FileCount:    info.FileCount,
		TotalSize:    info.TotalSize,
		LastModified: info.LastModified,
	})
}

// ListObjects 列出文件夹内的文件
func (s *S3Service) ListObjects(path, prefix string, limit int) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	path = normalizePath(path)
	if !strings.HasSuffix(path, "/") {
		path = path + "/"
	}

	// 完整前缀
	fullPrefix := path + prefix

	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:  aws.String(s.bucket),
		Prefix:  aws.String(fullPrefix),
		MaxKeys: aws.Int32(int32(limit)),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("列出文件失败: %v", err))
	}

	objects := []models.ObjectInfo{}
	for _, obj := range result.Contents {
		// 过滤掉文件夹标记（以 "/" 结尾的空对象）
		if strings.HasSuffix(*obj.Key, "/") {
			continue
		}

		// 生成临时访问 URL
		url, _ := s.generatePresignedURL(*obj.Key, 15)

		objects = append(objects, models.ObjectInfo{
			Key:          *obj.Key,
			Name:         filepath.Base(*obj.Key),
			Size:         *obj.Size,
			ContentType:  "",
			LastModified: obj.LastModified.Format(time.RFC3339),
			ETag:         strings.Trim(*obj.ETag, "\""),
			URL:          url,
		})
	}

	return common.SuccessWithData("获取成功", &models.ListObjectsResponse{
		Objects: objects,
		Total:   len(objects),
	})
}

// UploadFile 上传文件
func (s *S3Service) UploadFile(folder, filename string, content []byte) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	folder = normalizePath(folder)
	if !strings.HasSuffix(folder, "/") {
		folder = folder + "/"
	}

	// 构建完整 key
	key := folder + filename

	_, err := s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("上传失败: %v", err))
	}

	// 生成访问 URL
	url, _ := s.generatePresignedURL(key, 15)

	return common.SuccessWithData("上传成功", &models.UploadFileResponse{
		Key:  key,
		URL:  url,
		Size: int64(len(content)),
	})
}

// DeleteObject 删除文件
func (s *S3Service) DeleteObject(key string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("删除失败: %v", err))
	}

	return common.Success("删除成功")
}

// GetObjectInfo 获取单个文件详情
func (s *S3Service) GetObjectInfo(key string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	result, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("获取文件信息失败: %v", err))
	}

	// 生成临时访问 URL
	url, _ := s.generatePresignedURL(key, 15)

	info := &models.ObjectInfo{
		Key:          key,
		Name:         filepath.Base(key),
		Size:         *result.ContentLength,
		ContentType:  aws.ToString(result.ContentType),
		LastModified: result.LastModified.Format(time.RFC3339),
		ETag:         strings.Trim(aws.ToString(result.ETag), "\""),
		URL:          url,
	}

	return common.SuccessWithData("获取成功", info)
}

// DownloadObject 获取文件下载链接
func (s *S3Service) DownloadObject(key string, expireMinutes int) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	// 检查文件是否存在
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("文件不存在: %v", err))
	}

	// 默认 15 分钟过期
	if expireMinutes <= 0 {
		expireMinutes = 15
	}

	// 生成临时下载 URL
	url, err := s.generatePresignedURL(key, expireMinutes)
	if err != nil {
		return common.Error(fmt.Sprintf("生成下载链接失败: %v", err))
	}

	return common.SuccessWithData("获取成功", map[string]string{
		"key":       key,
		"url":       url,
		"expires_in": fmt.Sprintf("%dm", expireMinutes),
	})
}

// CopyObject 复制文件
func (s *S3Service) CopyObject(sourceKey, destKey string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	// 检查源文件是否存在
	_, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(sourceKey),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("源文件不存在: %v", err))
	}

	// 复制文件
	copySource := fmt.Sprintf("%s/%s", s.bucket, sourceKey)
	_, err = s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
		Bucket:     aws.String(s.bucket),
		Key:        aws.String(destKey),
		CopySource: aws.String(copySource),
	})
	if err != nil {
		return common.Error(fmt.Sprintf("复制失败: %v", err))
	}

	// 生成目标文件访问 URL
	url, _ := s.generatePresignedURL(destKey, 15)

	return common.SuccessWithData("复制成功", map[string]string{
		"source_key": sourceKey,
		"dest_key":   destKey,
		"url":        url,
	})
}

// BatchDeleteObjects 批量删除文件
func (s *S3Service) BatchDeleteObjects(keys []string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	if len(keys) == 0 {
		return common.Error("文件列表不能为空")
	}

	deleted := []string{}
	failed := []string{}

	// 分批删除（每次最多 1000 个）
	for i := 0; i < len(keys); i += 1000 {
		end := i + 1000
		if end > len(keys) {
			end = len(keys)
		}
		batch := keys[i:end]

		// 构建删除对象列表
		deleteObjects := []types.ObjectIdentifier{}
		for _, key := range batch {
			deleteObjects = append(deleteObjects, types.ObjectIdentifier{Key: aws.String(key)})
		}

		result, err := s.client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
			Bucket: aws.String(s.bucket),
			Delete: &types.Delete{
				Objects: deleteObjects,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			// 整批失败
			for _, key := range batch {
				failed = append(failed, key)
			}
			continue
		}

		// 处理部分失败
		if len(result.Errors) > 0 {
			failedMap := make(map[string]bool)
			for _, err := range result.Errors {
				failedMap[*err.Key] = true
				failed = append(failed, *err.Key)
			}
			for _, key := range batch {
				if !failedMap[key] {
					deleted = append(deleted, key)
				}
			}
		} else {
			deleted = append(deleted, batch...)
		}
	}

	return common.SuccessWithData("批量删除完成", &models.BatchDeleteResponse{
		Deleted: deleted,
		Failed:  failed,
		Total:   len(keys),
	})
}

// BatchCopyObjects 批量复制文件
func (s *S3Service) BatchCopyObjects(items []models.CopyItem) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	if len(items) == 0 {
		return common.Error("复制列表不能为空")
	}

	copied := []models.CopyResult{}
	failed := []string{}

	for _, item := range items {
		copySource := fmt.Sprintf("%s/%s", s.bucket, item.SourceKey)
		_, err := s.client.CopyObject(context.TODO(), &s3.CopyObjectInput{
			Bucket:     aws.String(s.bucket),
			Key:        aws.String(item.DestKey),
			CopySource: aws.String(copySource),
		})
		if err != nil {
			failed = append(failed, item.SourceKey)
			continue
		}

		url, _ := s.generatePresignedURL(item.DestKey, 15)
		copied = append(copied, models.CopyResult{
			SourceKey: item.SourceKey,
			DestKey:   item.DestKey,
			URL:       url,
		})
	}

	return common.SuccessWithData("批量复制完成", &models.BatchCopyResponse{
		Copied: copied,
		Failed: failed,
		Total:  len(items),
	})
}

// BatchGetObjectsInfo 批量获取文件信息
func (s *S3Service) BatchGetObjectsInfo(keys []string) *common.Response {
	if !utils.IsS3Enabled() {
		return common.Error("S3 服务未启用")
	}

	if len(keys) == 0 {
		return common.Error("文件列表不能为空")
	}

	objects := []models.ObjectInfo{}

	for _, key := range keys {
		result, err := s.client.HeadObject(context.TODO(), &s3.HeadObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(key),
		})
		if err != nil {
			continue // 跳过不存在的文件
		}

		url, _ := s.generatePresignedURL(key, 15)
		objects = append(objects, models.ObjectInfo{
			Key:          key,
			Name:         filepath.Base(key),
			Size:         *result.ContentLength,
			ContentType:  aws.ToString(result.ContentType),
			LastModified: result.LastModified.Format(time.RFC3339),
			ETag:         strings.Trim(aws.ToString(result.ETag), "\""),
			URL:          url,
		})
	}

	return common.SuccessWithData("获取成功", &models.BatchGetInfoResponse{
		Objects: objects,
		Total:   len(objects),
	})
}

// ===== 辅助方法 =====

// normalizePath 规范化路径
func normalizePath(path string) string {
	// 去掉开头的 "/"
	path = strings.TrimPrefix(path, "/")
	return path
}

// getFolderName 从路径获取文件夹名称
func getFolderName(path string) string {
	path = strings.TrimSuffix(path, "/")
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// folderExists 检查文件夹是否存在
func (s *S3Service) folderExists(path string) (bool, error) {
	result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(s.bucket),
		Prefix:    aws.String(path),
		MaxKeys:   aws.Int32(1),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return false, err
	}

	// 检查 CommonPrefixes 或 Contents
	for _, prefix := range result.CommonPrefixes {
		if *prefix.Prefix == path {
			return true, nil
		}
	}
	for _, obj := range result.Contents {
		if *obj.Key == path {
			return true, nil
		}
	}

	return false, nil
}

// listAllObjects 列出文件夹下所有对象
func (s *S3Service) listAllObjects(path string) ([]types.Object, error) {
	var allObjects []types.Object

	var continuationToken *string
	for {
		result, err := s.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			Prefix:            aws.String(path),
			ContinuationToken: continuationToken,
		})
		if err != nil {
			return nil, err
		}

		// 过滤掉文件夹标记
		for _, obj := range result.Contents {
			if !strings.HasSuffix(*obj.Key, "/") {
				allObjects = append(allObjects, obj)
			}
		}

		if result.IsTruncated == nil || !*result.IsTruncated {
			break
		}
		continuationToken = result.NextContinuationToken
	}

	return allObjects, nil
}

// getFolderStats 获取文件夹统计信息
func (s *S3Service) getFolderStats(path string) models.FolderInfo {
	objects, err := s.listAllObjects(path)
	if err != nil {
		return models.FolderInfo{}
	}

	var totalSize int64
	var lastModified time.Time

	for _, obj := range objects {
		totalSize += *obj.Size
		if obj.LastModified != nil && (*obj.LastModified).After(lastModified) {
			lastModified = *obj.LastModified
		}
	}

	return models.FolderInfo{
		FileCount:    int64(len(objects)),
		TotalSize:    totalSize,
		LastModified: lastModified.Format(time.RFC3339),
	}
}

// generatePresignedURL 生成临时访问 URL
func (s *S3Service) generatePresignedURL(key string, expireMinutes int) (string, error) {
	if s.presigner == nil {
		return "", fmt.Errorf("presigner 未初始化")
	}

	presignedURL, err := s.presigner.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(expireMinutes) * time.Minute
	})
	if err != nil {
		return "", err
	}

	return presignedURL.URL, nil
}