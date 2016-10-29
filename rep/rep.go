package rep

import (
	"github.com/couchbase/gocb"
	"fmt"
	"github.com/jeromedoucet/alienor-back/model"
)

const bucketName string = "alienor"

var (
	bucket *gocb.Bucket
)

// A repository provide some basic operations
// on data into data store
type Repository interface {
	Get(identifier string, document model.Document) (gocb.Cas, error)
	Insert(document model.Document) error
	Update(document model.Document, cas gocb.Cas) error
}

// todo close the bucket on exit
// prepare the repositories for requests.
// todo test me
func InitRepo(couchBaseAddr string, bucketPwd string) {
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
