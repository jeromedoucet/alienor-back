package rep_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/test"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/stretchr/testify/assert"
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
	assert.Nil(t, err)
	assert.True(t, exist)
}

func TestTeamExistWhenDoesNotTeamExist(t *testing.T) {
	// given
	test.Before()
	rep.InitRepo(test.CouchBaseAddr, "")

	// when
	exist, err := rep.TeamExist("A-Team", gocb.RequestPlus)

	// then
	assert.Nil(t, err)
	assert.False(t, exist)
}

