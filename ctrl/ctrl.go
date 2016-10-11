package ctrl

import "github.com/jeromedoucet/alienor-back/component"

var(
	rAdr string
	conn component.Connector = component.NewConnector();
)

func InitEndPoints(router component.Router, redisAddr string, secret string) {
	rAdr = redisAddr
	secr = []byte(secret)
	initAuthEndPoint(router)
	initUserEndPoint(router)
}
