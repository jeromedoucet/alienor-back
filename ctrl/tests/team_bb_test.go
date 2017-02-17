package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/component"
	"encoding/json"
	"bytes"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestTeamCreationSuccessFull(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name: "A-Team"}
	// prepare existing user
	usr := &model.User{Id: "leroy.jenkins", Type: model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL+"/team", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expected error to be nil")
	} else if res.StatusCode != 201 {
		t.Fatal("expected status code to equals 201")
	}

	// http res check
	var teamRes model.Team
	json.NewDecoder(res.Body).Decode(&teamRes)
	if teamRes.Name != teamCreationRequest.Name {
		t.Error("expect the name to be the same")
	}
	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(usr.Id)
	if len(actualUsr.Roles) != 1 {
		t.Error("expect to have only one role")
	} else if actualUsr.Roles[0].Value != model.ADMIN {
		t.Error("expect the role to be admin")
	} else if teamCreationRequest.Name != actualUsr.Roles[0].Team.Name {
		t.Error("expect the team name to eqauls")
	}
}

func TestTeamCreationWitoutBodyRequest(t *testing.T) {
	// given
	test.Before()
	// prepare existing user
	usr := &model.User{Id: "leroy.jenkins", Type: model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL+"/team", "POST", nil, token)

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 401 {
		t.Error("expect status code to equals 401")
	}

}

func TestTeamCreationWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name: "A-Team"}
	// prepare existing user
	usr := model.User{Id: "leroy.jenkins", Type: model.USER}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	res, err := test.DoReq(s.URL+"/team", "POST", bytes.NewBuffer(body))
	// then
	if err != nil {
		t.Error("expect err to be nil")
	} else if res.StatusCode != 401 {
		t.Error("expect status code to equals 401")
	}

	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(usr.Id)
	if len(actualUsr.Roles) != 0 {
		t.Error("expect user to have no role")
	}
}

// todo bench this test ! => first n1ql query => with a lot of data of course !
// todo check the error message too
func TestTeamCreationWhenTeamAlreadyExist(t *testing.T) {
	// given
	test.Before()
	teamCreationRequest := ctrl.TeamCreationReq{Name: "A-Team"}
	// prepare auth user
	leroy := model.User{Id: "leroy.jenkins", Type: model.USER}
	// prepare existing user with existing team
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + leroy.Id: leroy, "user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	body, _ := json.Marshal(teamCreationRequest)

	token := test.CreateToken(&leroy)
	res, err := test.DoReqWithToken(s.URL+"/team", "POST", bytes.NewBuffer(body), token)
	// then
	if err != nil {
		t.Error("expect err to be nil")
	} else if res.StatusCode != 409 {
		t.Error("expect status code to equals 409")
	}

	// db check -- the connected user should now be one admin of the
	actualUsr := test.GetUser(leroy.Id)
	if len(actualUsr.Roles) != 0 {
		t.Error("expect user to have no role")
	}
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
	usr := &model.User{Id: "leroy.jenkins", Type: model.USER}
	usr.Roles = []*model.Role{role}
	test.Populate(map[string]interface{}{"user:" + usr.Id: usr})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL+"/team", "GET", nil, token)

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 200 {
		t.Error("expect status code to equals")
	}
	var teamsRes []model.Team
	json.NewDecoder(res.Body).Decode(&teamsRes)

	if len(teamsRes) != 1 {
		t.Error("expect to have only one team")
	} else if teamsRes[0].Name != "the A-team" {
		t.Error("expect to team name to equals 'the A-team'")
	}
}

func TestTeamEnumerationWhenUserDoesntExist(t *testing.T) {
	// given
	test.Before()
	// prepare existing user with one team as Admin
	usr := &model.User{Id: "leroy.jenkins", Type: model.USER}

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	// when
	token := test.CreateToken(usr)
	res, err := test.DoReqWithToken(s.URL+"/team", "GET", nil, token)

	// then
	if err != nil {
		t.Error("expect err to be nil")
	} else if res.StatusCode != 404 {
		t.Error("expect status code to equals 404")
	}
}

func TestTeamEnumerationWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	// when
	res, err := test.DoReq(s.URL+"/team", "GET", nil)

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != 401 {
		t.Error("expect status code to equals 401")
	}
}
