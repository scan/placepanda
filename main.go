package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/h2non/bimg.v1"
	_ "github.com/joho/godotenv/autoload"
)

func corsMiddleware(h http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodOptions},
		Debug:            false,
	}).Handler(h)
}

func loggingMiddleware(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, h)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	w.Write([]byte("{\"version\":\"0.0.1\"}"))
}

func main() {
	bimg.Initialize()
	defer bimg.Shutdown()

	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/", versionHandler).Methods(http.MethodGet)
	router.HandleFunc("/panda/{width}/{height}", pandaHandler).Methods(http.MethodGet)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}).Methods(http.MethodGet)

	err := http.ListenAndServe("0.0.0.0:8080", handlers.RecoveryHandler()(corsMiddleware(router)))
	if err != nil {
		log.Println(err)
	}
}
