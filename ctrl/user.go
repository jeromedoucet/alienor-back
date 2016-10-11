package ctrl

import (
	"github.com/jeromedoucet/alienor-back/component"
	"net/http"
	"github.com/jeromedoucet/alienor-back/model"
	"encoding/json"
)

// user handler
func handleUser(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	var usr model.User
	dec.Decode(&usr) //todo handle error (with a bad formed json ?)

	c, _ := conn.Connect("tcp", rAdr) //todo handle error
	usrToSave, _ := json.Marshal(usr)
	c.Do("SET", usr.Identifier, string(usrToSave))
	w.Write(usrToSave)
}

func initUserEndPoint(router component.Router) {
	router.HandleFunc(UserEndPoint, handleUser)
}
