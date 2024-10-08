package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

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

func LogError(ctx context.Context, err error, failBlock string) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err == nil {
		err = errors.New("No ERROR object passed.")
	}
	logger.Error(
		err.Error(),
		"fail_block", failBlock,
		"trace_id", ctx.Value(ReqIdKey),
	)
}

func LASErrorResponse(r *http.Request, w http.ResponseWriter, status int, errorMessage string, err string, failBlock string) {
	GetLogEntry(r).Error(err, "fail_block", failBlock)
	ErrorResponse(w, status, errorMessage)
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
