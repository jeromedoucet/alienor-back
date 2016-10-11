package ctrl

import (
	"net/http"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"golang.org/x/crypto/bcrypt"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"errors"
	"time"
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
	c, cErr := conn.Connect("tcp", rAdr)
	if cErr != nil {
		w.WriteHeader(503)
		return
	}
	defer c.Close()
	exist, eErr := redis.Bool(c.Do("EXISTS", req.Login))
	if eErr != nil {
		w.WriteHeader(503)
		return
	}
	if exist {
		var user model.User
		bUser, _ := c.Do("GET", req.Login)
		json.Unmarshal(bUser.([]byte), &user)
		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Pwd)) == nil {
			token, jwtError := createJwtToken(user)
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

	} else {
		w.WriteHeader(404)
	}
}

// create the token used for the newly created session
func createJwtToken(usr model.User) (token string, err error) {
	// todo make the exp variable
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Identifier,
		"rls" : usr.Roles,
		"scp" : usr.Scope,
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
		if _, scp := claims["scp"]; !scp {
			err = errors.New("no scp claims")
			return
		}
		if _, rls := claims["rls"]; !rls {
			err = errors.New("no rls claims")
			return
		}
		usr.Roles = rolesFromClaim(claims)
		usr.Scope = scopeFromClaim(claims)
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

func scopeFromClaim(c jwt.MapClaims) []string {
	scp := c["scp"].([]interface{})
	scopes := make([]string, len(scp))
	for i, r := range scp {
		scopes[i] = r.(string)
	}
	return scopes
}

func rolesFromClaim(c jwt.MapClaims) []model.Role {
	rls := c["rls"].([]interface{})
	roles := make([]model.Role, len(rls))
	for i, r := range rls {
		roles[i] = model.Role(r.(string))
	}
	return roles
}