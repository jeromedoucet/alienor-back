package ctrl

import (
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/alienor-back/model"
	"testing"
	"net/http"
	"bytes"
	"net/http/httptest"
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/couchbase/gocb"
	"errors"
	"github.com/jeromedoucet/alienor-back/rep"
)

func TestIsLoggedWithSuccess(t *testing.T) {
	// given
	secret = []byte("some secret")
	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// when
	unMarshaledUsr, err := CheckToken(&r)
	assert.Nil(t, err)
	assert.Equal(t, usr.Identifier, unMarshaledUsr.Identifier)
}

func TestIsLoggedWithoutToken(t *testing.T) {
	// given
	r := http.Request{Header:http.Header{}}
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestIsLoggedWithBadToken(t *testing.T) {
	// given
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + "aBadTokenYouSee?")
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestIsLoggedWithExpiredToken(t *testing.T) {
	// given
	secret = []byte("some secret")

	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(-60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestIsLoggedWithoutBearerPrefix(t *testing.T) {
	// given
	secret = []byte("some secret")
	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", tokenString)
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestCheckUserCredentialBadRequestBody(t *testing.T) {
	// given
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBufferString("some string"))

	// when
	usr, err := checkUserCredential(req)

	// then
	assert.Nil(t, usr)
	assert.NotNil(t, err)
	assert.Equal(t, 400, err.httpCode)
	assert.Equal(t, "Error during decoding the authentication request body", err.errorMsg)
}

func TestCheckUserUnknownUser(t *testing.T) {
	// given
	userInReq := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	body, _ := json.Marshal(userInReq)
	req := httptest.NewRequest("POST", "http://127.0.0.1:8080", bytes.NewBuffer(body))
	defer func() {
		userRepository = new(rep.UserRepository) // reset userRepository
	}()
	userRepository = &utils.RepositoryHeader{DoGet:func(identifier string, entity interface{}) (gocb.Cas, error) {
		return 0, errors.New("some error")
	}}

	// when
	_, err := checkUserCredential(req)

	// then
	assert.NotNil(t, err)
	assert.Equal(t, 404, err.httpCode)
	assert.Equal(t, "Unknow User", err.errorMsg)
}

/* ################################################################################################################## */
/* ##############################################  BENCH  ########################################################### */
/* ################################################################################################################## */

func BenchmarkIsLogged(b *testing.B) {
	// given
	secret = []byte("some secret")
	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secret)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// bench
	for n := 0; n < b.N; n++ {
		CheckToken(&r)
	}
}

