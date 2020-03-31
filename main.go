package main

import (
	"almabase/Git-Ops/controller"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/getGitDetails", controller.GetGitDetails).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", r))

}
