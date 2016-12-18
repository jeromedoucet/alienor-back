package ctrl

import (
	"net/http"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/model"
	"fmt"
	"github.com/dgrijalva/jwt-go"
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
	secret []byte
)

// authentication handler
func handleAuth(w http.ResponseWriter, r *http.Request) {
	usr, err := checkUserCredential(r)
	if err != nil {
		writeError(w, err)
	} else {
		token, jwtError := createJwtToken(usr)
		if jwtError != nil {
			w.WriteHeader(500)
			return
		}
		writeSessionCookie(w, token)
		w.WriteHeader(200)
	}

}

// do perform req unmarshal, user fetch and password comparison
// as check during authentication
func checkUserCredential(r *http.Request) (usr *model.User, cError *ctrlError) {
	dec := json.NewDecoder(r.Body)
	var req AuthReq
	err := dec.Decode(&req)
	if err != nil {
		cError = &ctrlError{httpCode:400, errorMsg:"Error during decoding the authentication request body"}
		return
	}
	usr = model.NewUser()
	_, err = userRepository.Get(req.Login, usr)
	if err != nil {
		cError = &ctrlError{httpCode:404, errorMsg:"Unknow User"}
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(usr.Password), []byte(req.Pwd))
	if err != nil {
		cError = &ctrlError{httpCode:400, errorMsg:"Bad credentials"}
	}
	return
}

// create the token used for the newly created session
func createJwtToken(usr *model.User) (token string, err error) {
	// todo make the exp variable (and less than that !)
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usr.Id,
		"exp": time.Now().Add(1 * time.Minute).Unix(),
	})
	token, err = t.SignedString(secret)
	return
}

// write the session Cookie into the response
func writeSessionCookie(w http.ResponseWriter, token string)  {
	/*
	 * Assuming there is just one cookie for this app : the session one
	 * means that it must be cleared before set a new one.
	 */
	w.Header().Del("Set-Cookie")
	c := http.Cookie{}
	c.Name = "ALIENOR_SESS"
	c.HttpOnly = true
	c.Value = token
	http.SetCookie(w , &c)
}

// init the auth component by registering auth enpoint on router
func initAuthEndPoint(router component.Router) {
	router.HandleFunc(AUTH_ENDPOINT, handleAuth)
}

// this func will check the JWT token. If valid, a user is return
// an error otherwise.
func CheckToken(r *http.Request) (usr *model.User, err error) {
	var c *http.Cookie
	c, err = r.Cookie("ALIENOR_SESS")
	if err != nil {
		return
	}
	token, parsingError := jwt.Parse(string([]byte(c.Value)), keyFunc)

	if parsingError != nil {
		err = parsingError;
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		usr = new(model.User)
		usr.Id = claims["sub"].(string)
	} else {
		err = errors.New("invalid token or invalid claim type")
		return
	}
	return
}

func RefreshToken(w http.ResponseWriter, usr *model.User) (err error) {
	token, jwtError := createJwtToken(usr)
	if jwtError != nil {
		err = jwtError
		return
	}
	writeSessionCookie(w, token)
	return
}

// function which provide the secret
func keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}
	return secret, nil
}
