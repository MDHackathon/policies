package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/MDHackathon/policies/pkg"
	"github.com/minio/minio-go"
)

// TODO,md
const (
	bucket      string = "policies"
	access      string = "SCWWPMW5QBE5GHWNR02D"
	secret      string = "58b44520-dd5a-48c5-9727-d5d1889d23b5"
	endpoint    string = "beta.scalewaydata.com"
	letterBytes        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	client *minio.Client = nil
)

// generate a random name for the policies
// TODO get it from the json template
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Http handler for /policies
// will laod the data and stream it into
// the dedicated bucket
func pushPolicies(rw http.ResponseWriter, req *http.Request) {
	var (
		err error
	)

	if client, err = policies.GetS3Client(access, secret, endpoint, true); err != nil {
		panic("failed to initialized s3 client")
	}

	name := RandStringBytes(42)
	fmt.Println(name)
	// TODO find content length
	if _, err = client.PutObject(bucket, string(name), req.Body, -1, minio.PutObjectOptions{}); err != nil {
		panic("failed to initialized s3 client")
	}
}

// Open a http server on port 8080
// listen only the route /policies
// waiting for json policie to push on
// the dedicated s3 bucket
func main() {
	http.HandleFunc("/policies", pushPolicies)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
