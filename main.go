package main

import (
	"net/http"
	"os"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	models "github.com/edstef/goboyle/models"
)

var tokenAuth *jwtauth.JWTAuth

var mods *models.Models

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

	logLevel := 0
	mods = models.NewModels("postgres://postgres:@localhost:5432/sss?sslmode=disable", logLevel)

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

		registerProtectedJwtEndpoints(r)
		registerProtectedProfileEndpoints(r)
	})

	// Unprotected Routes
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("welcome"))
		})

		registerUnprotectedJwtEndpoints(r)
		registerUnprotectedProfileEndpoints(r)
	})

	http.ListenAndServe(":8080", r)
}
