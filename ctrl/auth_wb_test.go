package ctrl

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/alienor-back/model"
	"testing"
	"net/http"
	"bytes"
	"net/http/httptest"
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/couchbase/gocb"
	"errors"
	"github.com/jeromedoucet/alienor-back/rep"
	"golang.org/x/crypto/bcrypt"
)

func TestIsLoggedWithSuccess(t *testing.T) {
	// given
	t.Parallel()
	secret = []byte("some secret")
	usr := model.User{Id: "leroy.jenkins", Type: model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Id,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header: http.Header{}}
	// create session cookie
	c := http.Cookie{}
	c.Name = "ALIENOR_SESS"
	c.HttpOnly = true
	c.Value = tokenString
	r.AddCookie(&c)
	// when
	unMarshaledUsr, err := CheckToken(&r)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	} else if usr.Id != unMarshaledUsr.Id {
		t.Error("expect the unmarshalled usr id to equal the original one")
	}
}

func TestIsLoggedWithoutToken(t *testing.T) {
	// given
	t.Parallel()
	r := http.Request{Header: http.Header{}}
	// when
	_, err := CheckToken(&r)
	// then
	if err == nil {
		t.Error("expect the error not to be nil")
	}
}

func TestIsLoggedWithBadToken(t *testing.T) {
	// given
	t.Parallel()
	r := http.Request{Header: http.Header{}}
	c := http.Cookie{}
	// create session cookie
	c.Name = "ALIENOR_SESS"
	c.HttpOnly = true
	c.Value = "aBadTokenYouSee?"
	r.AddCookie(&c)
	// when
	_, err := CheckToken(&r)
	// then
	if err == nil {
		t.Error("expect the error not to be nil")
	}
}

func TestIsLoggedWithExpiredToken(t *testing.T) {
	// given
	t.Parallel()
	secret = []byte("some secret")

	usr := model.User{Id: "leroy.jenkins", Type: model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Id,
		"exp": time.Now().Add(-60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header: http.Header{}}
	// create session cookie
	c := http.Cookie{}
	c.Name = "ALIENOR_SESS"
	c.HttpOnly = true
	c.Value = tokenString
	r.AddCookie(&c)
	// when
	_, err := CheckToken(&r)
	// then
	if err == nil {
		t.Error("expect the error not to be nil")
	}
}

func TestCheckUserCredentialBadRequestBody(t *testing.T) {
	// given
	t.Parallel()
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBufferString("some string"))

	// when
	usr, err := checkUserCredential(req)

	// then
	if usr != nil {
		t.Error("expect the user to be nil")
	} else if err == nil {
		t.Error("expect the error not to be nil")
	} else if err.httpCode != 400 {
		t.Error("expect http return code to equals 400")
	} else if err.errorMsg != "Error during decoding the authentication request body" {
		t.Error("bad error message")
	}
}

func TestCheckUserUnknownUser(t *testing.T) {
	// given
	userInReq := model.User{Id: "leroy.jenkins", Type: model.USER}
	body, _ := json.Marshal(userInReq)
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(body))
	defer func() {
		userRepository = new(rep.UserRepository) // reset userRepository
	}()
	userRepository = &test.RepositoryHeader{DoGet: func(identifier string, document model.Document) (gocb.Cas, error) {
		return 0, errors.New("some error")
	}}

	// when
	_, err := checkUserCredential(req)

	// then
	if err == nil {
		t.Error("expect the error not to be nil")
	} else if err.httpCode != 404 {
		t.Error("expect http return code to equals 404")
	} else if err.errorMsg != "Unknow User" {
		t.Error("bad error message")
	}
}

func TestCheckUserBadPassword(t *testing.T) {
	// given
	userInReq := AuthReq{Login: "leroy.jenkins", Pwd: "wipe"}
	body, _ := json.Marshal(userInReq)
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(body))
	defer func() {
		userRepository = new(rep.UserRepository) // reset userRepository
	}()
	userRepository = &test.RepositoryHeader{DoGet: func(identifier string, document model.Document) (gocb.Cas, error) {
		userInRepo := document.(*model.User)
		userInRepo.Password = "roxxor"
		return 0, nil
	}}

	// when
	_, err := checkUserCredential(req)

	// then
	if err == nil {
		t.Error("expect the error not to be nil")
	} else if err.httpCode != 400 {
		t.Error("expect http return code to equals 404")
	} else if err.errorMsg != "Bad credentials" {
		t.Error("bad error message")
	}
}

func TestCheckUserSuccessFul(t *testing.T) {
	// given
	t.Parallel()
	userInReq := AuthReq{Login: "leroy.jenkins", Pwd: "wipe"}
	body, _ := json.Marshal(userInReq)
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(body))
	defer func() {
		userRepository = new(rep.UserRepository) // reset userRepository
	}()
	userRepository = &test.RepositoryHeader{DoGet: func(identifier string, document model.Document) (gocb.Cas, error) {
		pwd, _ := bcrypt.GenerateFromPassword([]byte("wipe"), bcrypt.DefaultCost)
		userInRepo := document.(*model.User)
		userInRepo.Password = string(pwd)
		userInRepo.Id = "leroy.jenkins"
		return 0, nil
	}}

	// when
	usr, err := checkUserCredential(req)

	// then
	if err != nil {
		t.Error("expect the error to be nil")
	} else if usr == nil {
		t.Error("expect the user not to be nil")
	} else if usr.Id != userInReq.Login {
		t.Error("expect the user id to eqauls the request login")
	}
}

func TestWriteSessionCookieNominalCase(t *testing.T) {
	// given
	t.Parallel()
	r := httptest.NewRecorder()
	token := "some-token"

	// when
	writeSessionCookie(r, token)

	// then
	resp := http.Response{Header: r.Header()}
	cookie := resp.Cookies()[0]

	if len(resp.Cookies()) != 1 {
		t.Error("expect to have only one cookie")
	} else if cookie.Name != "ALIENOR_SESS" {
		t.Error("bad cookie name")
	} else if cookie.Value != token {
		t.Error("bad cookie value")
	} else if !cookie.HttpOnly {
		t.Error("expect cookie in http only")
	}
}

func TestWriteSessionCookieWhenPreviousCookiesShouldRefreshIt(t *testing.T) {
	// given
	t.Parallel()
	r := httptest.NewRecorder()
	token := "some-token"

	c := http.Cookie{}
	c.Name = "ALIENOR_SESS"
	c.Value = "old-token"
	http.SetCookie(r, &c)

	// when
	writeSessionCookie(r, token)

	// then
	resp := http.Response{Header: r.Header()}
	cookie := resp.Cookies()[0]
	if len(resp.Cookies()) != 1 {
		t.Error("expect to have only one cookie")
	} else if cookie.Name != "ALIENOR_SESS" {
		t.Error("bad cookie name")
	} else if cookie.Value != token {
		t.Error("bad cookie value")
	} else if !cookie.HttpOnly {
		t.Error("expect cookie in http only")
	}
}

/* ################################################################################################################## */
/* ##############################################  BENCH  ########################################################### */
/* ################################################################################################################## */

// benchmark the check token function
func BenchmarkCheckToken(b *testing.B) {
	// given
	secret = []byte("some secret")
	usr := model.User{Id: "leroy.jenkins", Type: model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Id,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header: http.Header{}}
	r.Header.Set("Authorization", "bearer "+tokenString)
	// bench
	for n := 0; n < b.N; n++ {
		CheckToken(&r)
	}
}

// benchmark of the check user function (without db connection)
func BenchmarkChechUser(b *testing.B) {
	// given
	userInReq := AuthReq{Login: "leroy.jenkins", Pwd: "wipe"}
	body, _ := json.Marshal(userInReq)
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(body))
	defer func() {
		userRepository = new(rep.UserRepository) // reset userRepository
	}()
	userRepository = &test.RepositoryHeader{DoGet: func(identifier string, entity model.Document) (gocb.Cas, error) {
		pwd, _ := bcrypt.GenerateFromPassword([]byte("wipe"), bcrypt.DefaultCost)
		userInRepo := entity.(*model.User)
		userInRepo.Password = string(pwd)
		userInRepo.Id = "leroy.jenkins"
		return 0, nil
	}}

	// bench
	for n := 0; n < b.N; n++ {
		checkUserCredential(req)
	}
}
