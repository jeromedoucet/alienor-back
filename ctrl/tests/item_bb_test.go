package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
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
	utils.Before()
	illidan := utils.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	utils.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _:= json.Marshal(newItem)
	token := utils.CreateToken(illidan);
	// when
	res, err := utils.DoReqWithToken(s.URL + "/item?team-id=" + illidan.Roles[0].Team.Id, "POST", bytes.NewBuffer(body), token)

	// then
	assert.Nil(t, err)
	assert.Equal(t ,http.StatusCreated, res.StatusCode)
	// todo assert body response
	savedItem := utils.GetItem(newItem.Id, illidan.Roles[0].Team.Id)
	assert.NotNil(t, savedItem)
	assert.Equal(t, newItem.Id, savedItem.Id)
	assert.Equal(t, model.ITEM, savedItem.Type)
	assert.Equal(t, model.Newly , savedItem.State)
	assert.Equal(t, illidan.Roles[0].Team.Id , savedItem.TeamId)
}

func TestCreateItemWillFailedWhenNOtAuthenticated(t *testing.T) {
	// given
	utils.Before()
	illidan := utils.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	utils.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _:= json.Marshal(newItem)

	// when
	res, err := utils.DoReq(s.URL + "/item?team-id=" + illidan.Roles[0].Team.Id, "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t ,http.StatusUnauthorized, res.StatusCode)
}

// item id vide
// pas de team id
// json malforne
// not belonging team
// verification duplication cle par rapport a team