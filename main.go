package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var tokenAuth *jwtauth.JWTAuth

func init() {
	tokenAuth = jwtauth.New("HS256", []byte("secret"), nil) // TODO: Read secret from config
}

func main() {
	slogJSONHandler := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}.NewJSONHandler(os.Stdout)

	// Routes
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(NewStructuredLogger(slogJSONHandler))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// JWT Protected Routes
	r.Group(func(r chi.Router) {

		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/decode_jwt", func(w http.ResponseWriter, r *http.Request) {
			response := struct {
				UserId string `json:"user_id"`
			}{
				getUserIdFromContext(r.Context()),
			}

			SuccessResponse(w, response)
		})

		r.Post("/post_test", func(w http.ResponseWriter, r *http.Request) {
			params := struct {
				Param1 string `json:"param_1"`
				Param2 string `json:"param_2"`
			}{}

			err := getRequestBody(r, &params)
			if err != nil {
				ErrorResponse(w, http.StatusBadRequest, "")
				return
			}

			response := struct {
				UserId string `json:"user_id"`
				Param1 string `json:"param_1"`
				Param2 string `json:"param_2"`
			}{
				getUserIdFromContext(r.Context()),
				params.Param1,
				params.Param2,
			}

			SuccessResponse(w, response)
		})
	})


	// Unprotected Routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		r.Get("/wait", func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(1 * time.Second)
			LogEntrySetField(r, "wait", true)
			fmt.Println(middleware.GetReqID(r.Context()))
			GetLogEntry(r).Info("test")
			w.Write([]byte("hi"))

		})

		r.Get("/panic", func(w http.ResponseWriter, r *http.Request) {
			panic("oops")
		})

		r.Get("/add_fields", func(w http.ResponseWriter, r *http.Request) {
			LogEntrySetFields(r, map[string]interface{}{"foo": "bar", "bar": "foo"})
		})

		r.Get("/get_jwt/{id}", func(w http.ResponseWriter, r *http.Request) {
			_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": chi.URLParam(r, "id")})
			w.Write([]byte(tokenString))
		})
	})


	http.ListenAndServe(":8080", r)
}
