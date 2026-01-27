package s3

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	Client     *minio.Client
	BucketName string
}

func New(endpoint, accessKeyID, secretAccessKey, bucketName string) (*Minio, error) {
	const op = "object.minio.New"

	minioClient, err := minio.New(
		endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: false,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Minio{
		Client:     minioClient,
		BucketName: bucketName,
	}, nil
}

func (m *Minio) CreateBucketWithCheck(ctx context.Context, bucketName string) error {
	const op = "object.minio.CreateBucketWithCheck"
	err := m.Client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := m.Client.BucketExists(ctx, bucketName)
		if errBucketExists != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		if exists {
			return nil
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (m *Minio) PutObject(ctx context.Context, path string) error {
	const op = "object.minio.PutObject"

	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer file.Close()

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = m.Client.PutObject(
		ctx, m.BucketName, info.Name(), file, info.Size(), minio.PutObjectOptions{
			ContentType: "application/octet-stream",
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
