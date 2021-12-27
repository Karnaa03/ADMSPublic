package s3

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"git.solutions.im/XeroxAgriCensus/AgriTracking/conf"
)

type S3 struct {
	client *minio.Client
	region string
}

func New(config conf.Config) (s3 S3, err error) {
	client, err := minio.New(config.S3Config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.S3Config.AccessKey, config.S3Config.SecretKey, ""),
		Secure: config.S3Config.UseSSL,
	})
	if err != nil {
		return
	}
	s3 = S3{
		client: client,
		region: config.S3Config.Region,
	}
	return
}

func (s3 S3) Get(bucket, object string) (content []byte, err error) {
	ctx := context.Background()
	obj, err := s3.client.GetObject(ctx, bucket, object, minio.GetObjectOptions{})
	if err != nil {
		return
	}
	defer func() {
		err := obj.Close()
		if err != nil {
			return
		}
	}()

	stat, err := obj.Stat()
	if err != nil {
		return
	}
	content = make([]byte, stat.Size)
	_, err = obj.Read(content)
	if err != nil && err.Error() != io.EOF.Error() {
		return
	}
	return
}
