package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/rep"
)

// register and prepare the endpoints
func InitEndPoints(router component.Router, couchBaseAddr string, bucketPwd string, secret string) {
	rep.InitRepo(couchBaseAddr, bucketPwd)
	secr = []byte(secret)
	initAuthEndPoint(router)
	initUserEndPoint(router)
}
