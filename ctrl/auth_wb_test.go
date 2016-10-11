package ctrl

import (
	"github.com/stretchr/testify/assert"
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/jeromedoucet/alienor-back/model"
	"testing"
	"net/http"
)

func TestIsLoggedWithSuccess(t *testing.T) {
	// given
	secr = []byte("some secret")
	usr := model.User{Identifier:"leroy.jenkins", Roles:[]model.Role{model.TRANSLATOR}, Scope:[]string{"team 1"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"rls" : usr.Roles,
		"scp" : usr.Scope,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secr)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// when
	unMarshaledUsr, err := CheckToken(&r)
	assert.Nil(t, err)
	assert.Equal(t, usr, unMarshaledUsr)
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
	secr = []byte("some secret")
	usr := model.User{Identifier:"leroy.jenkins", Roles:[]model.Role{model.TRANSLATOR}, Scope:[]string{"team 1"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"rls" : usr.Roles,
		"scp" : usr.Scope,
		"exp": time.Now().Add(-60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secr)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestIsLoggedWithoutRequiredClaims(t *testing.T) {
	// given
	secr = []byte("some secret")
	usr := model.User{Identifier:"leroy.jenkins", Roles:[]model.Role{model.TRANSLATOR}, Scope:[]string{"team 1"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secr)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

func TestIsLoggedWithBearerPrefix(t *testing.T) {
	// given
	secr = []byte("some secret")
	usr := model.User{Identifier:"leroy.jenkins", Roles:[]model.Role{model.TRANSLATOR}, Scope:[]string{"team 1"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"rls" : usr.Roles,
		"scp" : usr.Scope,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secr)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", tokenString)
	// when
	_, err := CheckToken(&r)
	assert.NotNil(t, err)
}

/* ################################################################################################################## */
/* ##############################################  BENCH  ########################################################### */
/* ################################################################################################################## */

func BenchmarkIsLogged(b *testing.B) {
	// given
	secr = []byte("some secret")
	usr := model.User{Identifier:"leroy.jenkins", Roles:[]model.Role{model.TRANSLATOR}, Scope:[]string{"team 1"}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"rls" : usr.Roles,
		"scp" : usr.Scope,
		"exp": time.Now().Add(60 * time.Second).Unix(),
	})
	tokenString, _ := token.SignedString(secr)
	r := http.Request{Header:http.Header{}}
	r.Header.Set("Authorization", "bearer " + tokenString)
	// bench
	for n := 0; n < b.N; n++ {
		CheckToken(&r)
	}
}

