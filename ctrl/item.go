package ctrl

import (
	"encoding/json"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/jeromedoucet/alienor-back/route"
	"net/http"
	"net/url"
	"strings"
)

func initItemEndPoint(router component.Router) {
	router.HandleFunc(ITEMS_ENDPOINT, handleItems)
	router.HandleFunc(ITEM_ENDPOINT, handleItem)
}

/*
 * ************************************************************************************
 *	Requests on /team/:teamId/item endpoint
 * ************************************************************************************
 */

// this is the function that will handle
// request incoming on /team/:teamId/item resources
// will accept POST an get only
func handleItems(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	// the checks made here are :
	// => do we have a valid authentication ?
	// => is the request body correct ?
	// => is there an existing item ?
	// => is the authenticated user a member of the team we want to create the item into ?
	// if all the check are ok, the item is created and returned.
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: http.StatusUnauthorized})
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
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: http.StatusNotFound})
		return
	}
	teamId := route.SplitPath(r.URL.Path)[1]
	// we use the ame instance to check if an item already exist and if not to create
	// the new item
	_, err = itemRepository.Get(teamId, item.Id, &item)
	if err == nil {
		writeError(w, &ctrlError{errorMsg: "#ExistingItem", httpCode: http.StatusConflict})
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
		writeError(w, &ctrlError{errorMsg: "Unknown error during item creation", httpCode: http.StatusInternalServerError})
		return
	}
	writeJsonResponse(w, item, 201)
}

/*
 * ************************************************************************************
 *	Requests on /team/:teamId/item/:itemId endpoint
 * ************************************************************************************
 */

// this is the function that will handle
// request incoming on /team/:teamId/item resources
// will accept POST an get only
func handleItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		deleteItem(w, r)
	case http.MethodGet:
		getItem(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: http.StatusUnauthorized})
		return
	}
	var item model.Item
	cErr := unmarshallItem(r, &item)
	if cErr != nil {
		writeError(w, cErr)
		return
	}
	pPart := route.SplitPath(r.URL.Path)

	teamId, _ := url.PathUnescape(pPart[1])
	itemId, _ := url.PathUnescape(pPart[3])

	usr := model.NewUser()
	_, err = userRepository.Get(principal.Id, usr)
	// todo weird, but test me
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: http.StatusNotFound})
		return
	}
	cErr = checkTeamExistence(usr, teamId)
	if cErr != nil {
		writeError(w, cErr)
		return
	}
	err = itemRepository.Delete(teamId, itemId, &item)
	onDataSourceError(err, w)
}

func getItem(w http.ResponseWriter, r *http.Request) {
	// first check authentication token and retrieve the principal from the token
	principal, err := CheckToken(r)
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#NotAuthenticated", httpCode: http.StatusUnauthorized})
		return
	}
	pPart := route.SplitPath(r.URL.Path)
	teamId, _ := url.PathUnescape(pPart[1])
	itemId, _ := url.PathUnescape(pPart[3])

	usr := model.NewUser()
	_, err = userRepository.Get(principal.Id, usr)
	// todo weird, but test me
	if err != nil {
		writeError(w, &ctrlError{errorMsg: "#UnknownUser", httpCode: http.StatusNotFound})
		return
	}

	cErr := checkTeamExistence(usr, teamId)
	if cErr != nil {
		writeError(w, cErr)
		return
	}

	var item model.Item
	_, err = itemRepository.Get(teamId, itemId, &item)
	if ! onDataSourceError(err, w) {
		// no issue detected when getting the item
		body, _ := json.Marshal(&item)
		w.Write(body)
	}
}

func onDataSourceError(err error, w http.ResponseWriter) bool {
	if err != nil {
		if se, ok := err.(rep.DataSourceError); ok {
			if se.KeyNotFound() {
				writeError(w, &ctrlError{errorMsg: "#UnknownItem", httpCode: http.StatusNotFound})
				return true
			}
			if se.KeyExists() {
				writeError(w, &ctrlError{errorMsg: "#OutDatedVersion", httpCode: http.StatusConflict})
				return true
			}
		} else {
			// todo test that ?
			writeError(w, &ctrlError{errorMsg: err.Error(), httpCode: http.StatusInternalServerError})
			return true
		}
	}
	return false
}

/*
 * ************************************************************************************
 *	Commons resources
 * ************************************************************************************
 */

// unmarchall
func unmarshallItem(r *http.Request, item *model.Item) *ctrlError {
	err := json.NewDecoder(r.Body).Decode(item)
	if err != nil {
		return &ctrlError{errorMsg: "#BadRequestBody", httpCode: http.StatusBadRequest}

	}
	if strings.TrimSpace(item.Id) == "" {
		return &ctrlError{errorMsg: "#MissingItemIdentifier", httpCode: http.StatusBadRequest}

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
	return &ctrlError{errorMsg: "#UnknownTeam", httpCode: http.StatusNotFound}
}
