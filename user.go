package main

import (
	"net/http"

	models "github.com/edstef/goboyle/models"
	"github.com/go-chi/chi/v5"
)

func registerProtectedUserEndpoints(r chi.Router) {
	r.Get("/", getUser)
}

func registerUnprotectedUserEndpoints(r chi.Router) {
	r.Post("/user", createUser)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userId := getUserIdFromContext(r.Context())
	user, err := mods.GetUserById(userId)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "createUser - GetUserById")
		return
	}

	SuccessResponse(w, user)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Name string `json:"name"`
	}{}

	err := getRequestBody(r, &params)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "createUser - getRequestBody")
		return
	}

	user, err := mods.CreateUser(params.Name)
	if err != nil {
		LASErrorResponse(r, w, http.StatusBadRequest, "", err.Error(), "createUser - CreateUser")
		return
	}

	token := models.GenerateUUID()
	err = mods.UpsertUserToken(user.Id, token)
	if err != nil {
		LogError(r.Context(), err, "createUser - UpsertUserToken")
	}

	ret := struct {
		User  *models.User `json:"user"`
		Jwt   string       `json:"jwt"`
		Token string       `json:"token"`
	}{
		User:  user,
		Jwt:   createJWT(user.Id),
		Token: token,
	}

	SuccessResponse(w, &ret)
}
