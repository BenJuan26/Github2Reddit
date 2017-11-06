package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/BenJuan26/Github2Reddit/config"
	//"github.com/BenJuan26/Github2Reddit/reddit"
	"github.com/BenJuan26/Github2Reddit/github"
	"github.com/gorilla/mux"
)

var conf config.Config

func main() {
	conf = config.LoadConfig("config.json")

	router := mux.NewRouter()
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	fmt.Printf("Listening on port %d\n", conf.Port)
	err := http.ListenAndServe(":"+strconv.Itoa(conf.Port), router)
	if err != nil {
		panic(err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	release := github.ParseReleasePayload(r)

	fmt.Println(release.Name)
}
