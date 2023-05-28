package main

import (
	"context"

	models "github.com/edstef/goboyle/models"
	"github.com/uptrace/bun"
)

type Profile struct {
	Id         string `bun:",pk,notnull,type:uuid,default:uuid_generate_v4()"`
	Name       string `bun:"name,notnull"`
	PictureURL string `bun:"name,notnull,default:'/defaults/1'"`
	Theme      string `bun:"theme,notnull,default:'default_theme_1'"`
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	profile, err := mods.GetProfileById(id)
	if err != nil {
		// TODO: Log
		ErrorResponse(w, http.StatusBadRequest, "")
		return
	}

	SuccessResponse(w, profile)
}

func createProfile(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	err := getRequestBody(r, &params)
	if err != nil {
		// TODO: Log
		ErrorResponse(w, http.StatusBadRequest, "")
		return
	}

	mods.CreateProfile(params.Name)

	SuccessResponse(w, nil)
}
