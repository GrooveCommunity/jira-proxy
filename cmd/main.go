package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"os"

	"github.com/GrooveCommunity/proxy-jira/entity"
	"github.com/gorilla/mux"
)

type JiraStructure struct {
	ID        int
	Timestamp string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	log.Println("Port: ", os.Getenv("APP_PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var target interface{}

	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &target)

	log.Println(target)
}
