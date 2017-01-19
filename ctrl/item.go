package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/model"
)

type ItemCreationReq struct {
	Id string `json:"id"`
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	var req ItemCreationReq;
	json.NewDecoder(r.Body).Decode(&req)
	item := model.NewItem()
	item.Id = req.Id
	item.TeamId = r.URL.Query().Get("team-id")
	itemRepository.Insert(item)
	w.WriteHeader(201)
}

func initItemEndPoint(router component.Router) {
	authFilter := &AuthFilter{HandleBusiness:handleItem}
	router.HandleFunc(ITEM_ENDPOINT, authFilter.HandleAuth)
}
