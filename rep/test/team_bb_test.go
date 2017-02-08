package rep_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/couchbase/gocb"
)

func TestTeamExistWhenTeamExist(t *testing.T) {
	// given
	test.Before()
	illidan := test.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	test.Populate(map[string]interface{}{"user:" + illidan.Id: illidan})
	rep.InitRepo(test.CouchBaseAddr, "")

	// when
	exist, err := rep.TeamExist("A-Team", gocb.RequestPlus)

	// then
	if err != nil {
		t.Error("expect error not to be nil")
	} else if exist != true {
		t.Error("expect exist to be true")
	}
}

func TestTeamExistWhenDoesNotTeamExist(t *testing.T) {
	// given
	test.Before()
	rep.InitRepo(test.CouchBaseAddr, "")

	// when
	exist, err := rep.TeamExist("A-Team", gocb.RequestPlus)

	// then
	if err != nil {
		t.Error("expect err to be nil")
	} else if exist != false {
		t.Error("expect exist to be false")
	}
}

