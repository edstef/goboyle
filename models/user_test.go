package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	models "github.com/edstef/goboyle/models"
)

func TestGetUser(t *testing.T) {
	user, err := mods.CreateUser("edstef")
	assert.NotNil(t, user)
	assert.Nil(t, err)

	fetchedUser, fetchErr := mods.GetUserById(user.Id)
	assert.Equal(t, user.Id, fetchedUser.Id)
	assert.Nil(t, err)
	assert.Nil(t, fetchErr)
}

func TestUpsertAndValidateUserToken(t *testing.T) {
	user, err := mods.CreateUser("edstef")
	assert.NotNil(t, user)
	assert.Nil(t, err)

	token := models.GenerateUUID()
	err = mods.UpsertUserToken(user.Id, token)
	assert.Nil(t, err)

	err = mods.ValidateUserToken(user.Id, token)
	assert.Nil(t, err)

	err = mods.ValidateUserToken(user.Id, "bad-token")
	assert.NotNil(t, err)

	err = mods.UpsertUserToken(user.Id, "good-token")
	assert.Nil(t, err)

	err = mods.ValidateUserToken(user.Id, token)
	assert.NotNil(t, err)

	err = mods.ValidateUserToken(user.Id, "good-token")
	assert.Nil(t, err)
}
