package main

import (
	"encoding/json"
	"log"
	"net/http"

	"os"

	"github.com/GrooveCommunity/proxy-jira/entity"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/proxy-jira/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/proxy-jira/webhook", handleWebhook).Methods("POST")

	log.Println("Port: ", os.Getenv("APP_PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	log.Println(params)
}
