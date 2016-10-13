package ctrl_test

import (
	"net/http/httptest"
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"io"
	"crypto/tls"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
)

// shared data and config between ctrl bb test
var (
	rAddr string = "192.168.99.100:6379"
	secret string = "someSecret"
)

// exec registrator and start a tls server
func startHttp(registrator func(component.Router)) *httptest.Server {
	m := http.NewServeMux()
	registrator(m)
	return httptest.NewTLSServer(m)
}

// prepare a http request for testing. ca cert check is disable
func doReq(url string, verb string, reader io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(verb, url, reader)
	req.Header.Set("Content-Type", "application/json")
	// disable TSL cert chain because of httptest autosign cert
	cli := http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify:true}}}
	return cli.Do(req)
}

// clean the db
func clean(c redis.Conn) {
	c.Do("FLUSHDB")
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