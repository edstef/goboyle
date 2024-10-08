package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

func registerUnprotectedJwtEndpoints(r chi.Router) {
	r.Get("/{id}", createJWTHandler)
	r.Post("/", refreshJWTHandler)
}

func registerProtectedJwtEndpoints(r chi.Router) {
	r.Get("/decode_jwt", decodeJWTHandler)
}

func decodeJWTHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		UserId string `json:"user_id"`
	}{
		getUserIdFromContext(r.Context()),
	}

	SuccessResponse(w, response)
}

func createJWT(userId string) string {
	claims := map[string]interface{}{
		"user_id": userId,
	}

	jwtauth.SetExpiry(claims, time.Now().Add(time.Minute*15))
	_, tokenString, _ := tokenAuth.Encode(claims)
	return tokenString
}

func getUserIdFromJWT(token string) (string, error) {
	jwt, err := tokenAuth.Decode(token)
	if err != nil {
		return "", err
	}

	userId, success := jwt.Get("user_id")
	if success {
		return userId.(string), nil
	}
	return "", errors.New("No user ID found")
}

func createJWTHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(createJWT(chi.URLParam(r, "id"))))
}

func refreshJWTHandler(w http.ResponseWriter, r *http.Request) {
	params := struct {
		JWT       string `json:"jwt"`
		UserToken string `json:"user_token"`
	}{}

	err := getRequestBody(r, &params)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "refreshJWT - getRequestBody")
		return
	}

	userId, err := getUserIdFromJWT(params.JWT)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "refreshJWT - getUserIdFromJWT")
		return
	}

	err = mods.ValidateUserToken(userId, params.UserToken)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "refreshJWT - ValidateUserToken")
		return
	}

	SuccessResponse(w, createJWT(userId))
}
