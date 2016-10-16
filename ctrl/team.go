package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
)

// team creation means to be authenticated
// if not authenticated, then redirect to login
// check team identifier existence
// the creator is admin of the team
// if success, return 201

type TeamCreationReq struct {
	Name string `json:"name"`
}

// no need to return name, only identifier may be useful
type TeamCreationRes struct {
	Identifier string `json:"identifier"`
}

func handleTeam(w http.ResponseWriter, r *http.Request) {

}

// init the auth component by registering auth enpoint on router
func initTeamEndPoint(router component.Router) {
	router.HandleFunc(AUTH_ENDPOINT, handleAuth)
}
