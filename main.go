package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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
	http.ListenAndServe(":3333", r)
}
