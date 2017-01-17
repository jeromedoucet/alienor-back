package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/component"
)

func TestCreateItemNominal(t *testing.T) {
	// given
	utils.Before()
	//teamCreationRequest := ctrl.TeamCreationReq{Name:"A-Team"}
	// prepare auth user
	leroy := model.User{Id: "leroy.jenkins", Type:model.USER}
	// prepare existing user with existing team
	illidan := utils.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	utils.Populate(map[string]interface{}{"user:" + leroy.Id: leroy, "user:" + illidan.Id: illidan})

	s := utils.StartHttp(func(r component.Router) {ctrl.InitEndPoints(r, utils.CouchBaseAddr, "", utils.Secret)})
	defer s.Close()
}
