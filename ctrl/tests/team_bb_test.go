package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/component"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestTeamCreationSuccessFull(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare existing user
	usr := &model.User{Id: "leroy.jenkins", Type:model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL + "/team", "POST", bytes.NewBuffer(body), token)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 201, res.StatusCode)

	// http res check
	var teamRes model.Team
	json.NewDecoder(res.Body).Decode(&teamRes)
	assert.Equal(t, teamCreationRequest.Name, teamRes.Name)
	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(usr.Id)
	assert.Len(t, actualUsr.Roles, 1)
	assert.Equal(t, model.ADMIN, actualUsr.Roles[0].Value)
	assert.Equal(t, teamCreationRequest.Name, actualUsr.Roles[0].Team.Name)
}

func TestTeamCreationWitoutBodyRequest(t *testing.T) {
	// given
	test.Before()
	// prepare existing user
	usr := &model.User{Id: "leroy.jenkins", Type:model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL + "/team", "POST", nil, token)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)

}

func TestTeamCreationWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare existing user
	usr := model.User{Id: "leroy.jenkins", Type:model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	res, err := test.DoReq(s.URL + "/team", "POST", bytes.NewBuffer(body))
	// then
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)

	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(usr.Id)
	assert.Equal(t, 0, len(actualUsr.Roles))
}

// todo bench this test ! => first n1ql query => with a lot of data of course !
// todo check the error message too
func TestTeamCreationWhenTeamAlreadyExist(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare auth user
	leroy := model.User{Id: "leroy.jenkins", Type:model.USER}
	// prepare existing user with existing team
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + leroy.Id: leroy, "user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	token := test.CreateToken(&leroy)
	res, err := test.DoReqWithToken(s.URL + "/team", "POST", bytes.NewBuffer(body), token)
	// then
	assert.Nil(t, err)
	assert.Equal(t, 409, res.StatusCode)

	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(leroy.Id)
	assert.Equal(t, 0, len(actualUsr.Roles))
}

func TestTeamEnumerationSuccessFull(t *testing.T) {
	// given
	test.Before()
	// prepare existing user with one team as Admin
	team := model.NewTeam()
	team.Name = "the A-team"
	role := model.NewRole()
	role.Value = model.ADMIN
	role.Team = team
	usr := &model.User{Id: "leroy.jenkins", Type:model.USER}
	usr.Roles = []*model.Role{role}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL + "/team", "GET", nil, token)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	var teamsRes []model.Team
	json.NewDecoder(res.Body).Decode(&teamsRes)

	assert.Len(t, teamsRes, 1)
	assert.Equal(t, teamsRes[0].Name, "the A-team")
}

func TestTeamEnumerationWhenUserDoesntExist(t *testing.T) {
	// given
	test.Before()
	// prepare existing user with one team as Admin
	usr := &model.User{Id: "leroy.jenkins", Type:model.USER}

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL + "/team", "GET", nil, token)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 404, res.StatusCode)
}

func TestTeamEnumerationWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	// when
	res, err := test.DoReq(s.URL + "/team", "GET", nil)

	// then
	assert.Nil(t, err)
	assert.Equal(t, 401, res.StatusCode)
}