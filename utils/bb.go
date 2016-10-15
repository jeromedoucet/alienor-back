package utils

import (
	"net/http/httptest"
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"io"
	"crypto/tls"
	"github.com/couchbase/gocb"
	"github.com/jeromedoucet/alienor-back/model"
)

// todo check that this package is not in binary

// shared data and config between ctrl bb test
var (
	CouchBaseAddr string = "127.0.0.1"
	Cluster *gocb.Cluster
	Bucket *gocb.Bucket
	Secret string = "someSecret"
)

// initiate couchbase resources
func Before() {
	var err error
	Cluster, err = gocb.Connect("couchbase://" + CouchBaseAddr)
	if err != nil {
		panic(err)
	}
	Bucket, err = Cluster.OpenBucket("alienor", "")
	if err != nil {
		panic(err)
	}
}

func After() {
	Bucket.Close();

}

// exec registrator and start a tls server
func StartHttp(registrator func(component.Router)) *httptest.Server {
	m := http.NewServeMux()
	registrator(m)
	return httptest.NewTLSServer(m)
}

// prepare a http request for testing. ca cert check is disable
func DoReq(url string, verb string, reader io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(verb, url, reader)
	req.Header.Set("Content-Type", "application/json")
	// disable TSL cert chain because of httptest autosign cert
	cli := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}}
	return cli.Do(req)
}

// get one User
func GetUser(identifier string) (*model.User) {
	usr := model.NewUser()
	_, err := Bucket.Get(identifier, usr)
	if err != nil {
		panic(err)
	}
	return usr
}

// clean the db
func Clean(keys []string) {
	var items []gocb.BulkOp
	for _, k := range keys {
		items = append(items, &gocb.RemoveOp{Key: k})
	}
	doBulkOps(items)
}

func Populate(data map[string]interface{}) {
	var items []gocb.BulkOp
	for k, v := range data {
		items = append(items, &gocb.InsertOp{Key: k, Value: v})
	}
	doBulkOps(items)
}

func doBulkOps(items []gocb.BulkOp) {
	err := Bucket.Do(items)
	if err != nil {
		panic(err)
	}
}