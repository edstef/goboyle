package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProfile(t *testing.T) {
	profile, err := mods.CreateProfile("edstef")
	assert.NotNil(t, profile)
	assert.Nil(t, err)

	fetchedProfile, fetchErr := mods.GetProfileById(profile.Id)
	assert.Equal(t, profile.Id, fetchedProfile.Id)
	assert.Nil(t, err)
	assert.Nil(t, fetchErr)
}
