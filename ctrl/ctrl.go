package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/rep"
	"fmt"
	"net/http"
	"encoding/json"
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
func InitEndPoints(router component.Router, couchBaseAddr string, bucketPwd string, secret string) {
	rep.InitRepo(couchBaseAddr, bucketPwd)
	secr = []byte(secret)
	initAuthEndPoint(router)
	initUserEndPoint(router)
	initTeamEndPoint(router)
}

// todo test me !
func writeError(w http.ResponseWriter, err *ctrlError)  {
	errBody, _ := json.Marshal(ErrorBody{Msg:err.errorMsg})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.httpCode)
	fmt.Fprintf(w, "%s", errBody)
}
