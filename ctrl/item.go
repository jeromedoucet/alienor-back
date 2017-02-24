package ctrl

import (
	"net/http"
	"github.com/jeromedoucet/alienor-back/component"
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/route"
	"strings"
	"errors"
)

type ItemCreationReq struct {
	Id string `json:"id"`
}

func handleItem(w http.ResponseWriter, r *http.Request) {
	// todo check authenticated
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: 401})
		return
	}
	var req ItemCreationReq
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#BadRequestBody", httpCode: 400})
		return
	}
	if strings.TrimSpace(req.Id) == "" {
		writeError(w, &ctrlError{errorMsg: "#MissingItemIdentifier", httpCode: 400})
		return
	}
	usr := model.NewUser()
	_, err = userRepository.Get(principal.Id, usr)
	// todo weird, but test me
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: 404})
		return
	}
	item := model.NewItem()
	item.Id = req.Id
	item.TeamId = route.SplitPath(r.URL.Path)[1]
	err = checkTeamExistence(usr, item.TeamId)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: err.Error(), httpCode: 404})
		return
	}
	itemRepository.Insert(item)
	w.WriteHeader(201)
}

func checkTeamExistence(usr *model.User, teamId string) error {
	for _, r := range usr.Roles {
		if r.Team.Id == teamId {
			return nil
		}
	}
	return errors.New("#UnknownTeam")
}

func initItemEndPoint(router component.Router) {
	router.HandleFunc(ITEM_ENDPOINT, handleItem)
}
