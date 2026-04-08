package utils

import (
	"context"
	"log"
	"study-music-server-go/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var S3Presigner *s3.PresignClient

// InitS3Client 初始化 S3 客户端
func InitS3Client() error {
	cfg := config.AppConfig.AWS

	// 如果配置不完整，跳过初始化
	if cfg.Region == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		log.Println("AWS S3 配置不完整，跳过初始化")
		return nil
	}

	// 创建自定义 endpoint 解析器（用于兼容 S3 的服务）
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if cfg.Endpoint != "" {
			return aws.Endpoint{
				URL:           cfg.Endpoint,
				SigningRegion: region,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	// 加载 AWS 配置
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(cfg.Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		awsconfig.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return err
	}

	// 创建 S3 客户端
	S3Client = s3.NewFromConfig(awsCfg)

	// 创建预签名客户端（用于生成临时 URL）
	S3Presigner = s3.NewPresignClient(S3Client)

	log.Println("AWS S3 客户端初始化成功")
	return nil
}

// GetS3Bucket 获取 S3 bucket 名称
func GetS3Bucket() string {
	return config.AppConfig.AWS.Bucket
}

// IsS3Enabled 检查 S3 是否已启用
func IsS3Enabled() bool {
	return S3Client != nil
}