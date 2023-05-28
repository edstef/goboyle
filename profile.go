package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerProtectedProfileEndpoints(r chi.Router) {
	r.Get("/profile", getProfile)
}

func registerUnprotectedProfileEndpoints(r chi.Router) {
	r.Post("/profile", createProfile)
}

func getProfile(w http.ResponseWriter, r *http.Request) {
	profileId := getProfileIdFromContext(r.Context())
	profile, err := mods.GetProfileById(profileId)
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

	profile, err := mods.CreateProfile(params.Name)
	if err != nil {
		// TODO: Log
		ErrorResponse(w, http.StatusBadRequest, "")
		return
	}

	SuccessResponse(w, profile)
}
