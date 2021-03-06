package policies

import (
	"github.com/minio/minio-go"
)

// GetS3Client will initiate the s3 client
func GetS3Client(access string, secret string, endpoint string, ssl bool) (cli *minio.Client, err error) {
	cli, err = minio.NewV2(endpoint, access, secret, ssl)
	return
}
