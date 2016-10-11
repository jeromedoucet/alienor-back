package ctrl_test

import (
	"testing"
	"github.com/jeromedoucet/alienor-back/model"
	"github.com/jeromedoucet/alienor-back/component"
	"github.com/jeromedoucet/alienor-back/ctrl"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/garyburd/redigo/redis"
)

// nominal user creation test case
func TestUserCreationSuccessful(t *testing.T) {
	// given
	usr := model.User{Identifier:"leroy.jenkins",
		ForName:"Leroy",
		Name:"Jenkins",
		Email:"leroy.jenkins@wipe-guild.org",
		Password:[]byte("wipe"),
		}

	s := startHttp(func(r component.Router) {ctrl.InitEndPoints(r, rAddr, secret)})
	defer s.Close()
	body, _ := json.Marshal(usr)

	// when
	res, err := doReq(s.URL + "/user", "POST", bytes.NewBuffer(body))

	// then
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)

	// http res check
	var userRes model.User
	json.NewDecoder(res.Body).Decode(&userRes)
	assert.Equal(t, usr, userRes)

	// check db
	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	defer clean(c)
	var userDb model.User
	bUser, _ := c.Do("GET", usr.Identifier)
	json.Unmarshal(bUser.([]byte), &userDb)
	assert.Equal(t, usr, userDb)
}

// already used identifier
func TestUserCreationExistingIdentifier(t *testing.T) {

}

// when some mandatory fields are missing missing
func TestUserCreationMissingMandatoryField(t *testing.T) {

}


