package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
	"github.com/couchbase/gocb"
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
	principal, err := CheckToken(r)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	err = RefreshToken(w, principal)
	if err != nil {
		w.WriteHeader(401)
		return
	}
	switch r.Method {
	case http.MethodPost:
		createTeam(principal, w, r)
	case http.MethodGet:
		findTeam(principal, w)
	}
}

// do create a new team if possible
func createTeam(principal *model.User, w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var req TeamCreationReq
	err := dec.Decode(&req)
	if err != nil {
		writeError(w, &ctrlError{errorMsg:"#BadRequestBody", httpCode:401})
		return
	}

	ctrlErr := checkTeamExist(&req)
	if ctrlErr != nil {
		writeError(w, ctrlErr)
		return
	}

	var cas gocb.Cas
	usr := model.NewUser()
	cas, err = userRepository.Get(principal.Id, usr)
	if err != nil {
		// todo user nil ?? challenge me !
		// todo test me
		return
	}
	role := createNewRole(&req)
	if usr.Roles == nil {
		usr.Roles = []*model.Role{role}
	} else {
		usr.Roles = append(usr.Roles, role)
	}
	err = userRepository.Update(usr, cas)
	if err != nil {
		// todo test me
	}
	writeJsonResponse(w, role.Team, 201)
}

func findTeam(principal *model.User, w http.ResponseWriter) {
	_, err := userRepository.Get(principal.Id, principal)
	if err != nil {
		writeError(w, &ctrlError{errorMsg:"#UserNotFound", httpCode:404})
		return
	}
	teams := make([]*model.Team, len(principal.Roles), len(principal.Roles))
	for i, v := range principal.Roles {
		teams[i] = v.Team
	}
	writeJsonResponse(w, teams, 200)
}

// check if a team already exist with this name
func checkTeamExist(req *TeamCreationReq) *ctrlError {
	// todo check that the name is not empty
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
