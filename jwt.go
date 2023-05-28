package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func registerUnprotectedJwtEndpoints(r chi.Router) {
	r.Get("/get_jwt/{id}", createJWTHandler)
}

func registerProtectedJwtEndpoints(r chi.Router) {
	r.Get("/decode_jwt", decodeJWTHandler)
}

func decodeJWTHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		ProfileId string `json:"profile_id"`
	}{
		getProfileIdFromContext(r.Context()),
	}

	SuccessResponse(w, response)
}

func createJWTHandler(w http.ResponseWriter, r *http.Request) {
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"profile_id": chi.URLParam(r, "id")})
	w.Write([]byte(tokenString))
}
