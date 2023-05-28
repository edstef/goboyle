package main

import (
	"net/http"
	"os"
	"io/ioutil"

	"golang.org/x/exp/slog"
	"gopkg.in/yaml.v2"

	models "github.com/edstef/goboyle/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

var mods *models.Models

type Config struct {
	PgConnString string `yaml:"PG_CONNECTION_STRING"`
	Port         string `yaml:"PORT"`
	JwtSecret    string `yaml:"JWT_SECRET"`
}

func loadConfig(presets string) *Config {
	var c Config

	file, err := ioutil.ReadFile(presets)
	if err != nil {
		// logger.Fatal(err, fmt.Sprintf("Error reading from %s", presets))
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		// logger.Fatal(err, "Error unmarshalling yaml")
	}

	return &c
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

	conf := loadConfig("goboyle.yaml")
	logLevel := 0
	mods = models.NewModels(conf.PgConnString, logLevel)

	tokenAuth = jwtauth.New("HS256", []byte(conf.JwtSecret), nil)

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

	http.ListenAndServe(":" + conf.Port, r)
}
