package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/component"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestCreateItemNominal(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _:= json.Marshal(newItem)
	token := test.CreateToken(illidan);
	// when
	res, err := test.DoReqWithToken(s.URL + "/item?team-id=" + illidan.Roles[0].Team.Id, "POST", bytes.NewBuffer(body), token)

	// then
	assert.Nil(t, err)
	assert.Equal(t ,http.StatusCreated, res.StatusCode)
	// todo assert body response
	savedItem := test.GetItem(newItem.Id, illidan.Roles[0].Team.Id)
	assert.NotNil(t, savedItem)
	assert.Equal(t, newItem.Id, savedItem.Id)
	assert.Equal(t, model.ITEM, savedItem.Type)
	assert.Equal(t, model.Newly , savedItem.State)
	assert.Equal(t, illidan.Roles[0].Team.Id , savedItem.TeamId)
}

func TestCreateItemWillFailedWhenNOtAuthenticated(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret)})
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _:= json.Marshal(newItem)

	// when
	res, err := test.DoReq(s.URL + "/item?team-id=" + illidan.Roles[0].Team.Id, "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t ,http.StatusUnauthorized, res.StatusCode)
}

// item id vide
// pas de team id
// json malforne
// not belonging team
// verification duplication cle par rapport a team