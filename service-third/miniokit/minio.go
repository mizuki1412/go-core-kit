package miniokit

import (
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mizuki1412/go-core-kit/v2/class/exception"
	"github.com/mizuki1412/go-core-kit/v2/cli/configkey"
	"github.com/mizuki1412/go-core-kit/v2/service/configkit"
)

func NewClient() *minio.Client {
	minioClient, err := minio.New(configkit.GetString(configkey.MinioEndpoint), &minio.Options{
		Creds:  credentials.NewStaticV4(configkit.GetString(configkey.MinioAccessKey), configkit.GetString(configkey.MinioSecret), ""),
		Secure: false,
	})
	if err != nil {
		panic(exception.New(err.Error()))
	}
	return minioClient
}
