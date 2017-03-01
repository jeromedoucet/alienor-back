package rep

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/jeromedoucet/alienor-back/test"
)

func TestInsertItemSuccessFully(t *testing.T) {
	// given
	test.Before()
	item := model.NewItem()
	teamId := "teamId"
	item.Id = "#HelloWorld"
	item.Values = map[string]string{"fr_FR": ""}
	rep.InitRepo(test.CouchBaseAddr, "")
	itemRepository := new(rep.ItemRepository)
	// when
	err := itemRepository.Insert(teamId, item)

	// then
	if err != nil {
		t.Fatal("expect err to be nil")
	}
	savedItem := test.GetItem(item.Id, teamId)
	if savedItem.Id != item.Id {
		t.Fatal("expect id to be equals")
	} else if savedItem.Type != model.ITEM {
		t.Fatal("expect type to be ITEM")
	} else if savedItem.State != model.Newly {
		t.Fatal("expect state to be Newly")
	}
}

func TestInsertItemShouldFailedWhenEntityNotAnItem(t *testing.T) {
	// given
	test.Before()
	teamId := "teamId"
	rep.InitRepo(test.CouchBaseAddr, "")
	itemRepository := new(rep.ItemRepository)
	// when
	err := itemRepository.Insert(teamId, model.NewUser())

	// then
	if err == nil {
		t.Fatal("expect err not to be nil")
	} else if err.Error() != "can only insert item" {
		t.Fatal("wrong error type")
	}
}

func TestInsertItemShouldFailWhenNoType(t *testing.T) {
	// given
	test.Before()
	item := &model.Item{}
	teamId := "teamId"
	item.Id = "#HelloWorld"
	item.Values = map[string]string{"fr_FR": ""}
	rep.InitRepo(test.CouchBaseAddr, "")
	itemRepository := new(rep.ItemRepository)
	// when
	err := itemRepository.Insert(teamId, item)

	// then
	if err == nil {
		t.Fatal("expect err not to be nil")
	} else if err.Error() != "missing type in item" {
		t.Fatal("wrong error type")
	}
}

func TestInsertItemShouldFailWhenNoId(t *testing.T) {
	// given
	test.Before()
	item := model.NewItem()
	teamId := "teamId"
	item.Values = map[string]string{"fr_FR": ""}
	rep.InitRepo(test.CouchBaseAddr, "")
	itemRepository := new(rep.ItemRepository)
	// when
	err := itemRepository.Insert(teamId, item)

	// then
	if err == nil {
		t.Fatal("expect err not to be nil")
	} else if err.Error() != "missing id in item" {
		t.Fatal("wrong error type")
	}
}

func TestInsertItemShouldFailWhenBadStatus(t *testing.T) {
	// given
	test.Before()
	item := model.NewItem()
	teamId := "teamId"
	item.Id = "#HelloWorld"
	item.State = model.Complete
	item.Values = map[string]string{"fr_FR": ""}
	rep.InitRepo(test.CouchBaseAddr, "")
	itemRepository := new(rep.ItemRepository)
	// when
	err := itemRepository.Insert(teamId, item)

	// then
	if err == nil {
		t.Fatal("expect err not to be nil")
	} else if err.Error() != "bad item status" {
		t.Fatal("wrong error type")
	}
}

func TestGetItemSuccessFully(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	teamId := "team-id"
	existingItem := model.NewItem()
	existingItem.Id = itemId
	test.Populate(map[string]interface{}{"item:" + teamId + ":" + itemId: existingItem})
	itemRepository := new(rep.ItemRepository)
	newItem := &model.Item{}

	// when
	_, err := itemRepository.Get(teamId, itemId, newItem)

	// then
	if err != nil {
		t.Fatalf("expect error to be nil. Error is : %s", err.Error())
	} else if newItem.Id != existingItem.Id {
		t.Fatal("Id expected to be the same")
	} else if newItem.State != existingItem.State {
		t.Fatal("state expected to be the same")
	} else if newItem.Type != model.ITEM {
		t.Fatal("type expected to be ITEM")
	}
}

func TestGetItemShouldFailWhenDocumentIsNotAnItem(t *testing.T) {
	// given
	test.Before()
	itemId := "#HelloWorld"
	teamId := "team-id"
	existingItem := model.NewItem()
	existingItem.Id = itemId
	test.Populate(map[string]interface{}{"item:" + teamId + ":" + itemId: existingItem})
	itemRepository := new(rep.ItemRepository)
	newItem := &model.User{}

	// when
	_, err := itemRepository.Get(teamId, itemId, newItem)

	// then
	if err == nil {
		t.Fatal("expect error not to be nil")
	} else if err.Error() != "Cannot Get a non item entity" {
		t.Fatalf("bad error message : %s", err.Error())
	}
}
