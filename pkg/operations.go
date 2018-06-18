package policies

import (
	"github.com/minio/minio-go"
)

// Copy will load the data from the s3 input storage to the s3 output storage using the bucket
// input/output parameters
func Copy(in *minio.Client, out *minio.Client, ib string, ob string, objInfo minio.ObjectInfo) (err error) {
	var (
		stream *minio.Object
	)

	if stream, err = in.GetObject(ib, objInfo.Key, minio.GetObjectOptions{}); err != nil {
		return
	}
	// TODO check how use the n(int) return value
	if _, err = out.PutObject(ob, objInfo.Key, stream, objInfo.Size, minio.PutObjectOptions{}); err != nil {
		return
	}
	return

}

// Move will move the data from the s3 input storage to the s3 output storage
// using the input/output bucket.
func Move(in *minio.Client, out *minio.Client, ib string, ob string, objInfo minio.ObjectInfo) (err error) {

	if err = Copy(in, out, ib, ob, objInfo); err != nil {
		return
	}

	if err = in.RemoveObject(ib, objInfo.Key); err != nil {
		return
	}
	return
}
