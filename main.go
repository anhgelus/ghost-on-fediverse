package main

import (
	"github.com/anhgelus/ghost-on-fediverse/src"
	"github.com/gorilla/mux"
)

func main() {
	err := src.Connect()
	if err != nil {
		panic(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/", src.HandleWebhook)
}
