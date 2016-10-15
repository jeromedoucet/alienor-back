package ctrl

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
	"time"
	"github.com/jeromedoucet/alienor-back/rep"
)

type AuthReq struct {
	Login string `json:"login"`
	Pwd   string `json:"pwd"`
}

type AuthRes struct {
	Token string `json:"token"`
}

var (
	secr []byte
)

// authentication handler
func handleAuth(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req AuthReq
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(400)
		return
	}
	usr, eErr := rep.GetUser(req.Login)
	if eErr != nil {
		w.WriteHeader(404)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Pwd)) == nil {
		token, jwtError := createJwtToken(usr)
		if jwtError != nil {
			w.WriteHeader(500) //todo try to cover that (if possible)
			return
		}
		res, marshallError := json.Marshal(AuthRes{Token:token})
		if marshallError != nil {
			w.WriteHeader(500) //todo try to cover that (if possible)
			return
		}
		w.Write(res)
	} else {
		w.WriteHeader(400)
	}
}

// create the token used for the newly created session
func createJwtToken(usr *model.User) (token string, err error) {
	// todo make the exp variable
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"exp": time.Now().Add(20 * time.Minute).Unix(),
	})
	token, err = t.SignedString(secr)
	return
}

// init the auth component by registering auth enpoint on router
// setting redis addr and JWT HMAC secret for the run
func initAuthEndPoint(router component.Router) {
	router.HandleFunc(AuthEndpoint, handleAuth)
}

// this func will check the JWT token. If valid, a user is return
// an error otherwise.
func CheckToken(r *http.Request) (usr model.User, err error) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "bearer ") {
		err = errors.New("no valid authorization token")
		return
	}
	token, parsingError := jwt.Parse(string([]byte(auth)[7:]), keyFunc)

	if parsingError != nil {
		err = parsingError;
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usr.Identifier = claims["sub"].(string)
	} else {
		err = errors.New("invalid token or invalid claim type")
		return
	}
	return
}

// function which provide the secret
func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return secr, nil
}
