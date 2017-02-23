package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/component"
	"bytes"
	"encoding/json"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
)

func TestCreateItemNominal(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect error to be nil")
	} else if res.StatusCode != http.StatusCreated {
		t.Fatal("expect status code to equals 201")
	}
	// todo assert body response
	savedItem := test.GetItem(newItem.Id, illidan.Roles[0].Team.Id)
	if savedItem.Id != newItem.Id {
		t.Fatal("expect items id to be the same")
	} else if savedItem.Type != model.ITEM {
		t.Fatal("expect item type to be item")
	} else if savedItem.State != model.Newly {
		t.Fatal("expect item state to be new")
	} else if savedItem.TeamId != illidan.Roles[0].Team.Id {
		t.Fatal("bad team id")
	}
}

func TestCreateItemWillFailedWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)

	// when
	res, err := test.DoReq(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Error("expect error to be nil")
	} else if res.StatusCode != http.StatusUnauthorized {
		t.Error("expect status code to equals 401")
	}
}

func TestCreateItemWillFailedWhenItemIdEmpty(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := ctrl.ItemCreationReq{}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	} else if res.StatusCode != 400 {
		t.Fatal("expect http code to equals 400")
	}
	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#MissingItemIdentifier" {
		t.Fatal("wrong error msg")
	}

}

func TestCreateItemWillFailedWhenItemIdFilledWithBlank(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := ctrl.ItemCreationReq{Id:"           		\n\r"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	} else if res.StatusCode != 400 {
		t.Fatal("expect http code to equals 400")
	}
	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#MissingItemIdentifier" {
		t.Fatal("wrong error msg")
	}
}

// team id inexistant
// json malforne
// not belonging team
// verification duplication cle par rapport a team
