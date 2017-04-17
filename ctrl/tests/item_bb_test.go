package ctrl_test

import (
	"bytes"
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/test"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestCreateItemNominal(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect error to be nil")
	} else if res.StatusCode != http.StatusCreated {
		t.Fatalf("expect status code to equals 201 but was %d", res.StatusCode)
	}

	savedItem, cas := test.GetExistingItem(illidan.Roles[0].Team.Id, newItem.Id)
	if savedItem.Id != newItem.Id {
		t.Fatal("expect items id to be the same")
	} else if savedItem.Type != model.ITEM {
		t.Fatal("expect item type to be item")
	} else if savedItem.State != model.Newly {
		t.Fatal("expect item state to be new")
	}

	returnedItem := &model.Item{}
	json.NewDecoder(res.Body).Decode(returnedItem)
	if returnedItem.Id != newItem.Id {
		t.Fatal("expect items id to be the same")
	} else if returnedItem.Type != model.ITEM {
		t.Fatal("expect item type to be item")
	} else if returnedItem.State != model.Newly {
		t.Fatal("expect item state to be new")
	} else if uint64(cas) != returnedItem.Version {
		t.Fatal("expect item version to be equals")
	}
}

func TestCreateItemShouldReturn405WhenPut(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "PUT", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect error to be nil")
	} else if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expect status code to equals 405 but was %d", res.StatusCode)
	}
}

func TestCreateItemShouldReturn405WhenDelete(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "PUT", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect error to be nil")
	} else if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expect status code to equals 405 but was %d", res.StatusCode)
	}
}

func TestCreateItemWillFailedWhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
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

	newItem := model.Item{}
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

func TestCreateItemWillFailedWhenJsonMalformed(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer([]byte("something")), token)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	} else if res.StatusCode != 400 {
		t.Fatal("expect http code to equals 400")
	}

	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#BadRequestBody" {
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

	newItem := model.Item{Id: "           		\n\r"}
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

func TestCreateItemWillFailedWhenTeamNotFound(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/some-team/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	} else if res.StatusCode != 404 {
		t.Fatal("expect http code to equals 404")
	}

	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#UnknownTeam" {
		t.Fatal("wrong error msg")
	}
}

func TestCreateItemWillFailedWhenUsrNotFound(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: "#HelloWorld"}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	} else if res.StatusCode != 404 {
		t.Fatal("expect http code to equals 404")
	}

	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#UnknownUser" {
		t.Fatal("wrong error msg")
	}
}

func TestCreateShouldFailedWhenItemAlreadyExist(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	existingItem := model.NewItem()
	existingItem.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: existingItem})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: itemId}
	body, _ := json.Marshal(newItem)
	token := test.CreateToken(illidan)
	// when
	res, err := test.DoReqWithToken(s.URL+"/team/"+illidan.Roles[0].Team.Id+"/item", "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatal("expect error to be nil")
	} else if res.StatusCode != 409 {
		t.Fatal("expect status code to equals 409")
	}

	var errBody ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&errBody)
	if errBody.Msg != "#ExistingItem" {
		t.Fatal("wrong error msg")
	}
}

func TestDeleteItemShouldReturn401WhenNotAuthenticated(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	newItem := model.Item{Id: itemId}
	body, _ := json.Marshal(newItem)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReq(s.URL+path, "DELETE", bytes.NewBuffer(body))

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 401 {
		t.Fatalf("expect http code to be 401, but is %d", res.StatusCode)
	}
}

func TestDeleteItemShouldSucceed(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "DELETE", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 200 {
		t.Fatalf("expect http code to be 200, but is %d", res.StatusCode)
	}
	_, err = test.GetItem(illidan.Roles[0].Team.Id, itemId)
	if err == nil {
		t.Fatal("expect item to be deleted")
	}
}

func TestDeleteItemShouldFailWhenPost(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "POST", bytes.NewBuffer(body), token)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 405 {
		t.Fatalf("expect http code to be 405, but is %d", res.StatusCode)
	}
}

func TestDeleteItemShouldFailWhenInvalidBody(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(itemId)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "DELETE", strings.NewReader("something"), token)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 400 {
		t.Fatalf("expect http code to be 400, but is %d", res.StatusCode)
	}
}

func TestDeleteItemShouldFailWhenItemDoesNotExist(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	item := model.Item{Id: "some-item"}

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(item)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "DELETE", bytes.NewBuffer(body), token)

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 404 {
		t.Fatalf("expect http code to be 404, but is %d", res.StatusCode)
	} else if resB.Msg != "#UnknownItem" {
		t.Fatalf("expect error message to be '##UnknownItem', but is %s", resB.Msg)
	}
}

func TestDeleteItemShouldFailWhenItemVersionIsOutdated(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	item := model.NewItem()
	item.Id = "item-id"
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + item.Id: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, item.Id)
	itemToDelete.Version = itemToDelete.Version - 1

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "DELETE", bytes.NewBuffer(body), token)

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != http.StatusConflict {
		t.Fatalf("expect http code to be 409, but is %d", res.StatusCode)
	} else if resB.Msg != "#OutDatedVersion" {
		t.Fatalf("expect error message to be '#OutDatedVersion', but is %s", resB.Msg)
	}
}

func TestDeleteItemShouldFailWhenTeamDoesNotExist(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/toto/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "DELETE", bytes.NewBuffer(body), token)

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != 404 {
		t.Fatalf("expect http code to be 404, but is %d", res.StatusCode)
	} else if resB.Msg != "#UnknownTeam" {
		t.Fatalf("expect error message to be '#UnknownTeam', but is %s", resB.Msg)
	}
}

func TestGetItemShouldSucceed(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Values = map[string]string{"fr_FR": "quelque chose", "en_EN": "something"}
	item.State = model.Incomplete
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "GET", bytes.NewBuffer(body), token)

	var i model.Item
	json.NewDecoder(res.Body).Decode(&i)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != http.StatusOK {
		t.Fatalf("expect the status to be 200, but got %d", res.StatusCode)
	} else if i.Id != itemId {
		t.Fatalf("expect the identifier to be the same, but got %s and %s", i.Id, itemId)
	} else if len(i.Values) != 2 {
		t.Fatalf("expect the value length to be 2, but got %d", len(i.Values))
	} else if i.Values["fr_FR"] != "quelque chose" {
		t.Fatalf("expect the french translation to be 'quelque chose', but got %s", i.Values["fr_FR"])
	} else if i.Values["en_EN"] != "something" {
		t.Fatalf("expect the english translation to be 'something', but got %s", i.Values["fr_FR"])
	} else if i.State != model.Incomplete {
		t.Fatalf("expect the state to be Incomplete, but is %s", i.State)
	}
}

func TestGetItemShouldFailWhenNotAuth(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Values = map[string]string{"fr_FR": "quelque chose", "en_EN": "something"}
	item.State = model.Incomplete
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:" + illidan.Roles[0].Team.Id + ":" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReq(s.URL+path, "GET", bytes.NewBuffer(body))

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expect the status to be 401, but got %d", res.StatusCode)
	} else if resB.Msg != "#NotAuthenticated" {
		t.Fatalf("expect the error msg to be '#NotAuthenticated', but got %s", resB.Msg)
	}
}

func TestGetItemShouldFailWhenUnknownTeam(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	item := model.NewItem()
	item.Values = map[string]string{"fr_FR": "quelque chose", "en_EN": "something"}
	item.State = model.Incomplete
	item.Id = itemId
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan, "item:The-new-team:" + itemId: item})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape("The-new-team") + "/item/" + url.PathEscape(item.Id)

	// when
	res, err := test.DoReqWithToken(s.URL+path, "GET", bytes.NewBuffer(body), token)

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expect the status to be 404, but got %d", res.StatusCode)
	} else if resB.Msg != "#UnknownTeam" {
		t.Fatalf("expect the error msg to be '#UnknownTeam', but got %s", resB.Msg)
	}
}

func TestGetItemShouldFailWhenUnknownItem(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})

	itemToDelete, _ := test.GetItem(illidan.Roles[0].Team.Id, itemId)

	s := test.StartHttp(func(r component.Router) { ctrl.InitEndPoints(r, test.CouchBaseAddr, "", test.Secret) })
	defer s.Close()

	token := test.CreateToken(illidan)
	body, _ := json.Marshal(itemToDelete)
	path := "/team/" + url.PathEscape(illidan.Roles[0].Team.Id) + "/item/" + url.PathEscape("#HelloWorld")

	// when
	res, err := test.DoReqWithToken(s.URL+path, "GET", bytes.NewBuffer(body), token)

	var resB ctrl.ErrorBody
	json.NewDecoder(res.Body).Decode(&resB)
	if err != nil {
		t.Fatalf("expect error to be nil, but is %s", err.Error())
	} else if res.StatusCode != http.StatusNotFound {
		t.Fatalf("expect the status to be 404, but got %d", res.StatusCode)
	} else if resB.Msg != "#UnknownItem" {
		t.Fatalf("expect the error msg to be '#UnknownItem', but got %s", resB.Msg)
	}
}

// todo list item avec pagination
// todo maj item
// todo suppr item

// verification duplication cle par rapport a team
