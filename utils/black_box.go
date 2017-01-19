package utils

import (
	"net/http/httptest"
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"io"
	"crypto/tls"
	"github.com/couchbase/gocb"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// todo check that this package is not in binary

// shared data and config between ctrl bb test
var (
	CouchBaseAddr string = "127.0.0.1"
	Cluster *gocb.Cluster
	Secret string = "someSecret"
)

// initiate couchbase resources
func Before() {
	var err error
	var Bucket *gocb.Bucket
	Cluster, err = gocb.Connect("couchbase://" + CouchBaseAddr)
	if err != nil {
		panic(err)
	}
	Bucket, err = Cluster.OpenBucket("alienor", "")
	if err != nil {
		panic(err)
	}
	defer Bucket.Close()
	bManager := Bucket.Manager("alienor", "")
	err = bManager.CreatePrimaryIndex("aliaIndex", true, false)
	if err != nil {
		panic(err)
	}
	query := gocb.NewN1qlQuery("DELETE FROM alienor")
	query.Consistency(gocb.RequestPlus)
	_, err = Bucket.ExecuteN1qlQuery(query, []interface{}{})
	if err != nil {
		panic(err)
	}
}

// exec registrator and start a tls server
func StartHttp(registrator func(component.Router)) *httptest.Server {
	m := http.NewServeMux()
	registrator(m)
	return httptest.NewTLSServer(m)
}

// prepare a http request for testing. ca cert check is disable
func DoReq(url string, verb string, reader io.Reader) (*http.Response, error) {
	return DoReqWithToken(url, verb, reader, "")
}

func DoReqWithToken(url string, verb string, reader io.Reader, token string) (*http.Response, error) {
	req, _ := http.NewRequest(verb, url, reader)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		c := http.Cookie{}
		c.Name = "ALIENOR_SESS"
		c.HttpOnly = true
		c.Value = token
		req.AddCookie(&c)
	}
	// disable TSL cert chain because of httptest autosign cert
	cli := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}}
	return cli.Do(req)
}

// get one User
func GetUser(identifier string) (*model.User) {
	Bucket, _ := Cluster.OpenBucket("alienor", "")
	defer Bucket.Close()
	usr := model.NewUser()
	_, err := Bucket.Get("user:" + identifier, usr)
	if err != nil {
		panic(err)
	}
	return usr
}

// get one User
func GetItem(itemId, teamId string) (*model.Item) {
	Bucket, _ := Cluster.OpenBucket("alienor", "")
	defer Bucket.Close()
	item := model.NewItem()
	_, err := Bucket.Get(string(model.ITEM) + ":" + teamId + ":" + itemId, item)
	if err != nil {
		panic(err)
	}
	return item
}

func CreateToken(usr *model.User) (token string) {
	var err error
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Id,
		"exp": time.Now().Add(20 * time.Minute).Unix(),
	})
	token, err = t.SignedString([]byte(Secret))
	if err != nil {
		panic(err.Error())
	}
	return
}

func PrepareUserWithTeam(teamName string, identifier string) *model.User {
	team := model.NewTeam()
	team.Name = teamName
	role := model.NewRole()
	role.Team = team
	user := model.NewUser()
	user.Id = identifier
	user.Roles = []*model.Role{role}
	return user
}

func Populate(data map[string]interface{}) {
	bucket, _ := Cluster.OpenBucket("alienor", "")
	defer bucket.Close()
	var items []gocb.BulkOp
	for k, v := range data {
		items = append(items, &gocb.UpsertOp{Key: k, Value: v})
	}
	doBulkOps(items, bucket)
}

func doBulkOps(items []gocb.BulkOp, bucket *gocb.Bucket) {
	err := bucket.Do(items)
	if err != nil {
		panic(err)
	}
}