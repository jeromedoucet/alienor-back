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

type DataSourceError interface {
	KeyNotFound() bool
	KeyExists() bool
}

// todo interface for repositories are maybe not useful. Should we remove them ?

// A repository that fit for operation on root Documents :
// no need to provide another identifier than the document one: it is unique
type RootEntityRepository interface {
	Get(id string, doc model.Document) (gocb.Cas, error)
	Insert(doc model.Document) error
	Update(doc model.Document, cas gocb.Cas) error
}

// A repository that fit for operation on child Documents, stored
// a complex identifier
type ChildEntityRepository interface {
	Get(parentId, id string, doc model.Document) (gocb.Cas, error)
	Insert(parentId string, doc model.Document) error
	Update(parentId string, doc model.Document, cas gocb.Cas) error
	Delete(parentId, id string, doc model.Document) error
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
