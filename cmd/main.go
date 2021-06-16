package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GrooveCommunity/proxy-jira/entity"
	"github.com/GrooveCommunity/proxy-jira/internal"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	log.Println("Port: ", os.Getenv("APP_PORT"))

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	log.Println("Acionando healthy")
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var jiraRequest entity.JiraRequest

	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &jiraRequest)

	log.Println("Enviando para a fila do dispatcher")

	go internal.ForwardDispatcher("monitoria-groovetech", jiraRequest)
}
