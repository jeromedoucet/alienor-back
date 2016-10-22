package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/component"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestTeamCreationSuccessFull(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare existing user
	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	utils.Populate(map[string]interface{}{"user:" + usr.Identifier: usr})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	token := utils.CreateToken(&usr)
	res, err := utils.DoReqWithToken(s.URL + "/team", "POST", bytes.NewBuffer(body), token)
	// then
	assert.Nil(t, err)
	assert.Equal(t, 201, res.StatusCode)

	// http res check
	var teamRes model.Team
	json.NewDecoder(res.Body).Decode(&teamRes)
	assert.Equal(t, teamCreationRequest.Name, teamRes.Name)
	// db check -- the connected user should now be one admin of the
	actualUsr := utils.GetUser(usr.Identifier)
	assert.Equal(t, 1, len(actualUsr.Roles))
	assert.Equal(t, model.ADMIN, actualUsr.Roles[0].Value)
	assert.Equal(t, teamCreationRequest.Name, actualUsr.Roles[0].Team.Name)
}

func TestTeamCreationWhenNotAuthenticated(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare existing user
	usr := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	utils.Populate(map[string]interface{}{"user:" + usr.Identifier: usr})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	res, err := utils.DoReq(s.URL + "/team", "POST", bytes.NewBuffer(body))
	// then
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)

	// db check -- the connected user should now be one admin of the
	actualUsr := utils.GetUser(usr.Identifier)
	assert.Equal(t, 0, len(actualUsr.Roles))
}

// todo bench this test ! => first n1ql query
// todo check the error message too
func TestTeamCreationWhenTeamAlreadyExist(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare auth user
	leroy := model.User{Identifier: "leroy.jenkins", Type:model.USER}
	// prepare existing user with existing team
	illidan := utils.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	utils.Populate(map[string]interface{}{"user:" + leroy.Identifier: leroy, "user:" + illidan.Identifier: illidan})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	token := utils.CreateToken(&leroy)
	res, err := utils.DoReqWithToken(s.URL + "/team", "POST", bytes.NewBuffer(body), token)
	// then
	assert.Nil(t, err)
	assert.Equal(t, 409, res.StatusCode)

	// db check -- the connected user should now be one admin of the
	actualUsr := utils.GetUser(leroy.Identifier)
	assert.Equal(t, 0, len(actualUsr.Roles))
}