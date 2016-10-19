package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"github.com/couchbase/gocb"
	"fmt"
)

// team creation means to be authenticated
// if not authenticated, then redirect to login
// check team identifier existence
// the creator is admin of the team
// if success, return 201

type TeamCreationReq struct {
	Name string `json:"name"`
}

func handleTeam(w http.ResponseWriter, r *http.Request) {
	var usr *model.User
	principal, err := CheckToken(r)
	if err != nil {
		// todo test me
		fmt.Println(err.Error())
		return
	}
	dec := json.NewDecoder(r.Body)
	var req TeamCreationReq
	err = dec.Decode(&req)
	// todo check team existance
	var cas gocb.Cas
	usr, cas = rep.GetUser(principal.Identifier)
	if usr == nil {
		// todo test me
		fmt.Println(err.Error())
		return
	}
	role := createNewRole(&req)
	if usr.Roles == nil {
		usr.Roles = []*model.Role{role}
	} else {
		usr.Roles = append(usr.Roles, role)
	}
	err = rep.UpdateUser(usr, cas)
	if err != nil {
		// todo handle me
	}
	newTeam, _ := json.Marshal(role.Team)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", newTeam)
}

func createNewRole(req *TeamCreationReq) *model.Role {
	role := model.NewRole()
	role.Value = model.ADMIN
	role.Team = model.NewTeam()
	role.Team.Name = req.Name
	return role
}

// init the auth component by registering auth enpoint on router
func initTeamEndPoint(router component.Router) {
	router.HandleFunc(TEAM_ENDPOINT, handleTeam)
}
