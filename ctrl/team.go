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
		w.WriteHeader(401)
		return
	}
	dec := json.NewDecoder(r.Body)
	var req TeamCreationReq
	err = dec.Decode(&req)

	ctrlErr := checkTeamExist(&req)
	if ctrlErr != nil {
		writeError(w, ctrlErr)
		return
	}

	var cas gocb.Cas
	usr, cas = rep.GetUser(principal.Identifier)
	if usr == nil { // todo user nil ?? challenge me !
		// todo test me
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
		// todo test me
	}
	newTeam, _ := json.Marshal(role.Team)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", newTeam)
}

func checkTeamExist(req *TeamCreationReq) *ctrlError {
	exist, err := rep.TeamExist(req.Name, gocb.RequestPlus)
	if err != nil {
		return &ctrlError{httpCode:503, errorMsg:"Error during fetching data from the data store"}
	}
	if exist {
		return &ctrlError{httpCode:409, errorMsg:"Error during creating the team : already exist"}
	}
	return nil
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
