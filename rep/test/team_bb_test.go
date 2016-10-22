package rep_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/utils"
	"github.com/jeromedoucet/alienor-back/rep"
	"github.com/stretchr/testify/assert"
	"github.com/couchbase/gocb"
)

func TestTeamExistWhenTeamExist(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	illidan := utils.PrepareUserWithTeam("A-Team", "illidan.stormrage")
	utils.Populate(map[string]interface{}{"user:" + illidan.Identifier: illidan})
	rep.InitRepo(utils.CouchBaseAddr, "")

	// when
	exist, err := rep.TeamExist("A-Team", gocb.RequestPlus)

	// then
	assert.Nil(t, err)
	assert.True(t, exist)
}

func TestTeamExistWhenDoesNotTeamExist(t *testing.T) {
	// given
	utils.Before()
	defer utils.Clean()
	rep.InitRepo(utils.CouchBaseAddr, "")

	// when
	exist, err := rep.TeamExist("A-Team", gocb.RequestPlus)

	// then
	assert.Nil(t, err)
	assert.False(t, exist)
}

