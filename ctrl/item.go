package ctrl

import (
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/route"
	"net/http"
	"net/url"
	"strings"
)

// this is the function that will handle
// request incoming on /team/:teamId/item resources
// will accept POST an get only
func handleItems(w http.ResponseWriter, r *http.Request) {
	// the checks made here are :
	// => do we have a valid authentication ?
	// => is the request body correct ?
	// => is there an existing item ?
	// => is the authenticated user a member of the team we want to create the item into ?
	// if all the check are ok, the item is created and returned.
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: 401})
		return
	}
	item := model.Item{}
	cErr := unmarshallItem(r, &item)
	if cErr != nil {
		writeError(w, cErr)
		return
	}
	usr := model.NewUser()
	_, err = userRepository.Get(principal.Id, usr)
	// todo weird, but test me
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: 404})
		return
	}
	teamId := route.SplitPath(r.URL.Path)[1]
	// we use the ame instance to check if an item already exist and if not to create
	// the new item
	_, err = itemRepository.Get(teamId, item.Id, &item)
	if err == nil {
		writeError(w, &ctrlError{errorMsg: "#ExistingItem", httpCode: 409})
		return
	}
	cErr = checkTeamExistence(usr, teamId)
	if cErr != nil {
		writeError(w, cErr)
		return
	}
	err = itemRepository.Insert(teamId, &item)
	if err != nil {
		// todo need some log here
		writeError(w, &ctrlError{errorMsg: "Unknown error during item creation", httpCode: 500})
		return
	}
	writeJsonResponse(w, item, 201)
}

// this is the function that will handle
// request incoming on /team/:teamId/item resources
// will accept POST an get only
func handleItem(w http.ResponseWriter, r *http.Request) {
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: 401})
		return
	}
	var item model.Item
	cErr := unmarshallItem(r, &item)
	if cErr != nil {
		writeError(w, cErr)
	}
	pPart := route.SplitPath(r.URL.Path)

	teamId, _ := url.PathUnescape(pPart[1]) // todo test that unit test in function
	itemId, _ := url.PathUnescape(pPart[3]) // todo test that unit test in function

	usr := model.NewUser()
	_, err = userRepository.Get(principal.Id, usr)
	// todo weird, but test me
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: 404})
		return
	}
	cErr = checkTeamExistence(usr, teamId)
	if cErr != nil {
		writeError(w, cErr)
		return
	}
	itemRepository.Delete(teamId, itemId, &item)

}

// unmarchall
func unmarshallItem(r *http.Request, item *model.Item) *ctrlError {
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		return &ctrlError{errorMsg: "#BadRequestBody", httpCode: 400}

	}
	if strings.TrimSpace(item.Id) == "" {
		return &ctrlError{errorMsg: "#MissingItemIdentifier", httpCode: 400}

	}
	item.Type = model.ITEM
	item.State = model.Newly
	return nil
}

func checkTeamExistence(usr *model.User, teamId string) *ctrlError {
	for _, r := range usr.Roles {
		if r.Team.Id == teamId {
			return nil
		}
	}

	return &ctrlError{errorMsg: "#UnknownTeam", httpCode: 404}
}

func initItemEndPoint(router component.Router) {
	router.HandleFunc(ITEMS_ENDPOINT, handleItems)
	router.HandleFunc(ITEM_ENDPOINT, handleItem)
}
