package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/rep"
	"fmt"
	"net/http"
	"encoding/json"
)

var (
	userRepository rep.Repository = new(rep.UserRepository)
)

type ctrlError struct {
	httpCode int
	errorMsg string
}

type ErrorBody struct {
	Msg string `json:"msg"`
}

func (e *ctrlError) Error() string {
	return fmt.Sprintf("Error during crtl treatement : %s ", e.errorMsg)
}

// register and prepare the endpoints
func InitEndPoints(router component.Router, couchBaseAddr string, bucketPwd string, s string) {
	rep.InitRepo(couchBaseAddr, bucketPwd)
	secret = []byte(s)
	initAuthEndPoint(router)
	initUserEndPoint(router)
	initTeamEndPoint(router)
}

// write the error directly on the given response writer
func writeError(w http.ResponseWriter, err *ctrlError) {
	writeJsonResponse(w, ErrorBody{Msg:err.errorMsg}, err.httpCode)
}

// write an arbitrary response on the writer with the desired http code
// todo handle the marshall err !
func writeJsonResponse(w http.ResponseWriter, data interface{}, code int) {
	body, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, "%s", body)
}
