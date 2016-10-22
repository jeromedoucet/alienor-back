package rep

import (
	"github.com/couchbase/gocb"
	"fmt"
)

const bucketName string = "alienor"

var (
	bucket *gocb.Bucket
)

// prepare the repositories for requests.
// todo test me
func InitRepo(couchBaseAddr string, bucketPwd string)  {
	// todo create the bucket if needed !
	cluster, err := gocb.Connect("couchbase://" + couchBaseAddr)
	if err != nil {
		fmt.Println("ERRROR CONNECTING TO CLUSTER:", err) // todo test me
		panic(err)
	}
	// open it one time, it's thread-safe
	bucket, err = cluster.OpenBucket(bucketName, bucketPwd)
	if err != nil {
		fmt.Println("ERRROR OPENING BUCKET:", err) // todo test me
		panic(err)
	}
}
