package ctrl_test

import (
	"encoding/json"
	"testing"
	"net/http"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"net/http/httptest"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"github.com/jeromedoucet/alienor-back/model/team"
	"github.com/stretchr/testify/assert"
	"bytes"
	"crypto/tls"
	"io"
)

var rAddr string = "192.168.99.100:6379";

func TestHandleAuthSuccess(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := team.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(rAddr)
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:pwd})

	// when
	res, err := doReq(s.URL + "/login", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	// todo check jwt
}

func TestHandleRedisConError(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := team.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp("192.168.99.100:1234")
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:pwd})

	// when
	res, err := doReq(s.URL + "/login", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 503, res.StatusCode)
}

func TestHandleBadPassword(t *testing.T) {
	// given
	pwd := "wipe"
	login := "leroy.jenkins"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := team.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(rAddr)
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:login, Pwd:"roxx"})

	// when
	res, err := doReq(s.URL + "/login", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 400, res.StatusCode)
}

func TestHandleUnknownUser(t *testing.T) {
	// given
	pwd := "lichking"
	login := "arthas.menethil"
	hPwd, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost);
	usr := team.User{Identifier:login, Password:hPwd}

	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	populate(c, map[string]interface{}{usr.Identifier: usr})

	s := startHttp(rAddr)
	defer s.Close()
	body, _ := json.Marshal(ctrl.AuthReq{Login:"leroy.jenkins", Pwd:"test"})

	// when
	res, err := doReq(s.URL + "/login", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func startHttp(redisUrl string) *httptest.Server {
	m := http.NewServeMux()
	ctrl.InitAuth(m, redisUrl)
	return httptest.NewTLSServer(m)
}

func populate(c redis.Conn, data map[string]interface{}) {
	buf := make([]interface{}, len(data) * 2)
	for k, v := range data {
		buf = append(buf, k)
		val, _ := json.Marshal(v)
		buf = append(buf, string(val))
	}
	c.Do("MSET", buf...)

}

func clean(c redis.Conn) {
	c.Do("FLUSHDB")
}

func doReq(url string, reader io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest("POST", url, reader)
	req.Header.Set("Content-Type", "application/json")
	// disable TSL cert chain because of httptest autosign cert
	cli := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}}
	return cli.Do(req)
}