package storage

import (
	"bytes"
	"context"

	"github.com/Resolution-hash/shop_bot/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinIOClinent(cfg *config.Config) (*minio.Client, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: false,
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func MinIOPutPhoto(filename string, data []byte) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	minioClient, err := InitMinIOClinent(cfg)
	if err != nil {
		return err
	}
	_, err = minioClient.PutObject(context.Background(), cfg.Backet, filename, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{ContentType: "image/jpeg"})
	if err != nil {
		return err
	}
	return nil
}

func MinIORemovePhoto(filename string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}

	minioClient, err := InitMinIOClinent(cfg)
	if err != nil {
		return err
	}
	err = minioClient.RemoveObject(context.Background(), cfg.Backet, filename, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

func MinIOGetPhoto(filename string) (*minio.Object, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	minioClient, err := InitMinIOClinent(cfg)
	if err != nil {
		return nil, err
	}
	object, err := minioClient.GetObject(context.Background(), cfg.Backet, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
