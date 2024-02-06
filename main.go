package main

import (
	"github.com/anhgelus/ghost-on-fediverse/src"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	err = src.ConnectMastodon()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", src.HandleWebhook)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	src.LogInfo("Starting server...")
	src.Crash(srv.ListenAndServe())
}
