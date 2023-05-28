package main

import (
	"net/http"
	"encoding/json"
	"context"

	"github.com/go-chi/jwtauth/v5"

)

func getUserIdFromContext(ctx context.Context) string {
	_, claims, _ := jwtauth.FromContext(ctx)
	return claims["user_id"].(string)
}

func getRequestBody(r *http.Request, holder interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(holder)
	return err
}

func ErrorResponse(w http.ResponseWriter, status int, errorMessage string) {
	GenericResponse(w, errorMessage, status)
}

func SuccessResponse(w http.ResponseWriter, data interface{}, status ...int) {
	returnCode := 200
	if len(status) > 0 {
		returnCode = status[0]
	}

	GenericResponse(w, data, returnCode)
}

func GenericResponse(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
