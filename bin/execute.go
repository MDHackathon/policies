package main

import (
	"fmt"
	"github.com/MDHackathon/policies/pkg"
	"github.com/minio/minio-go"
	"io"
)

const (
	bucket   string = "policies"
	access   string = "SCWWPMW5QBE5GHWNR02D"
	secret   string = "58b44520-dd5a-48c5-9727-d5d1889d23b5"
	endpoint string = "beta.scalewaydata.com"
)

// This will load the policies from the s3 bucket policies define as const
// and apply the different policies
// the policies rules.
//
//	S3 storage:
//		- In  -> the storage source
//		- Out -> the storage destination
//
// Rules:
// 		- match(name) -> allow a reggex on the key
//		- match(date) -> allow to operate an a set of object define by date
//						 Format "2013-06-18T02:08:03+00:00"
//
// Operation:
//		- Copy -> Copy an object from a s3 storage to an other ( actif <-> passif )
//		- Move -> Moce an object from a s3 storage to an other (  hot  <->  cold  )
//
func main() {

	var (
		cli *minio.Client
		obj *minio.Object
		b   []byte
		err error
	)

	fmt.Println("it work !")
	if cli, err = policies.GetS3Client(access, secret, endpoint, true); err != nil {
		panic("failed to initialized s3 client")
	}

	doneCh := make(chan struct{})
	defer close(doneCh)

	objectCh := cli.ListObjects(bucket, "", true, doneCh)
	for objInfo := range objectCh {
		if obj, err = cli.GetObject(bucket, objInfo.Key, minio.GetObjectOptions{}); err != nil {
			fmt.Println("Failed to get the object", objInfo.Key)
			continue
		}
		for {
			b = make([]byte, objInfo.Size)
			_, err := obj.Read(b)
			if err == io.EOF {
				break
			}
		}
		policies.Execute(b)
	}
}
