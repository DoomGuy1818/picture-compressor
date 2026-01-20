package s3

import (
	"fmt"

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
		return nil, fmt.Errorf("%w: %s", err, op)
	}

	return &Minio{
		Client:     minioClient,
		BucketName: bucketName,
	}, nil
}

func (m *Minio) PutObject() error {
	panic("implement me")
}
