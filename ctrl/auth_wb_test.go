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
	usr, err :=checkUserCredential(req)

	// then
	assert.Nil(t, usr)
	assert.NotNil(t, err)
	assert.Equal(t, 400, err.httpCode)
	assert.Equal(t, "Error during decoding the authentication request body", err.errorMsg)
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

