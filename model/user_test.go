package model_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/stretchr/testify/assert"
	"github.com/satori/go.uuid"
)

func TestNewTeam(t *testing.T) {
	// when
	team := model.NewTeam()
	id, err := uuid.FromString(team.Id)
	assert.Nil(t, err)
	assert.NotEqual(t, "", id)

}

// benchmark the new team creation. go.uuid lib is used for that purposed
func BenchmarkNewTeam(b *testing.B) {
	// bench
	for n := 0; n < b.N; n++ {
		model.NewTeam()
	}
}