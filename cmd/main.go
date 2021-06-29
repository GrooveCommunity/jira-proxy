package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
	"github.com/GrooveCommunity/proxy-jira/internal"
	"github.com/gorilla/mux"
)

var (
	projectID, topicDispatcher, topicMetrics string
)

type Result struct {
	Teste string
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/healthy", handleValidateHealthy).Methods("GET")
	router.HandleFunc("/webhook", handleWebhook).Methods("POST")

	projectID = os.Getenv("PROJECT_ID")
	topicDispatcher = os.Getenv("TOPIC_ID_DISPATCHER")
	topicMetrics = os.Getenv("TOPIC_ID_METRICS")

	if projectID == "" || topicDispatcher == "" || topicMetrics == "" || os.Getenv("TOKEN_DISPATCHER_PAYGO_DISCORD") == "" || os.Getenv("CHANNEL_ID_DISCORD") == "" {
		log.Fatal("Nem todas as variáveis de ambiente requeridas foram fornecidas. ")
	}

	internal.GetNocUsers()

	log.Fatal(http.ListenAndServe(":"+os.Getenv("APP_PORT"), router))
}

func handleValidateHealthy(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(entity.Healthy{Status: "Success!"})
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	var request interface{}
	var jiraRequest entity.JiraRequest

	body, _ := ioutil.ReadAll(r.Body)

	json.Unmarshal(body, &jiraRequest)
	json.Unmarshal(body, &request)

	log.Println(string(body))
	log.Println("=========================================\n\n\n")

	customFields := internal.UnmarchallMapCustomField(request.(map[string]interface{}))

	jiraRequest.Issue.Fields.CustomFields = customFields

	go internal.ForwardIssue(jiraRequest, body, projectID, topicDispatcher, topicMetrics)
}
