package internal

import (
	"encoding/json"
	"reflect"

	"log"
	"strings"
	"time"

	"github.com/GrooveCommunity/glib-cloud-storage/gcp"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

type customFields map[string]interface{}

const YELLOW = 16705372

var nocUsers entity.NocUsers

func ForwardIssue(jiraRequest entity.JiraRequest, body []byte, projectID, topicDispatcher, topicMetrics string) {
	jiraEvent := entity.JiraEvent{
		EventUser: jiraRequest.User.Name,
		DateTime:  time.Now().Format(time.RFC3339),
		EventName: jiraRequest.EventName,
	}

	var jiraTransitions []entity.JiraTransition

	for _, change := range jiraRequest.ChangeLog.Changes {
		jiraTransitions = append(jiraTransitions, entity.JiraTransition{LastState: change.From, CurrentState: change.To})
	}

	jiraIssue := entity.JiraIssue{
		Event:        jiraEvent,
		CustomFields: jiraRequest.Issue.Fields.CustomFields,
		Transitions:  jiraTransitions,
		IssueID:      jiraRequest.Issue.ID,
		Key:          jiraRequest.Issue.Key,
		Assignee:     jiraRequest.Issue.Fields.Assignee.User,
		Creator:      jiraRequest.Issue.Fields.Creator.Name,
		Reporter:     jiraRequest.Issue.Fields.Reporter.Name,
		ChangeDate:   jiraRequest.Issue.Fields.ChangeDate,
		Priority:     jiraRequest.Issue.Fields.Priority.Name,
		Status:       jiraRequest.Issue.Fields.Status.Name,
		Type:         jiraRequest.Issue.Fields.IssueType.Name,
		CreatedDate:  jiraRequest.Issue.Fields.Created,
		UpdatedDate:  jiraRequest.Issue.Fields.Updated,
		Summary:      jiraRequest.Issue.Fields.Summary,
		Description:  jiraRequest.Issue.Fields.Description,
		Project:      jiraRequest.Issue.Fields.Project.Name,
		Attachment:   jiraRequest.Issue.Fields.Attachment,
	}

	payload, errPayLoad := json.Marshal(jiraIssue)

	if errPayLoad != nil {
		log.Fatal(entity.ResponseError{
			Message:    "Erro na conversão do payload para JSON",
			StatusCode: 500,
			Error:      errPayLoad,
		})
	}

	go validateIssueDispatcher(jiraRequest, projectID, topicDispatcher, payload)

	go PublicMessage(projectID, topicMetrics, payload)
}

func validateIssueDispatcher(jiraRequest entity.JiraRequest, projectID, topicName string, payload []byte) {
	msg := "Começou um novo ciclo de SLA!\n"

	for _, item := range jiraRequest.Issue.Fields.CustomFields {
		//customfield_10646 é o campo Squads
		if item.CustomID == "customfield_10366" {

			if validateForward(jiraRequest, item.Name, item.Value) {
				PublicMessage(projectID, topicName, payload)
				SendMessageToChannel(
					"https://paygo.atlassian.net/browse/"+jiraRequest.Issue.Key,
					jiraRequest.Issue.Key,
					msg+"\nPrioridade: "+jiraRequest.Issue.Fields.Priority.Name+"\nSLA: "+getSLA(jiraRequest.Issue.Fields.Priority.Name)+"\n\n\n",
					YELLOW)
			}

			break
		}
	}
}

func validateSLAUser(user string) bool {
	for _, nocUser := range nocUsers.JiraUsers {
		if nocUser.Name == user {
			return true
		}
	}

	return false
}

func UnmarchallMapCustomField(dataMap map[string]interface{}) []entity.CustomField {

	var customFields []entity.CustomField

	for key, item := range dataMap {
		if (key == "issue" || key == "fields") && reflect.TypeOf(item).Kind() == reflect.Map {
			customFields = UnmarchallMapCustomField(item.(map[string]interface{}))
		}

		if strings.HasPrefix(key, "customfield") && item != nil && item != "" && item != "{}" {
			jsonItem, errJsonItem := json.Marshal(item)
			if errJsonItem != nil {
				panic(errJsonItem)
			}

			var customField entity.CustomField

			json.Unmarshal(jsonItem, &customField)

			customField.CustomID = key

			customFields = append(customFields, customField)
		}
	}

	return customFields
}

func getSLA(priority string) string {
	if priority == "Altissima" {
		return "00:15:00"
	} else if priority == "Alta" {
		return "02:00:00"
	} else if priority == "Media" {
		return "08:00:00"
	} else if priority == "Baixa" {
		return "48:00:00"
	}

	return ""
}

func validateForward(jiraRequest entity.JiraRequest, customFieldName, customFieldValue string) bool {
	if jiraRequest.EventName == "jira:issue_updated" && jiraRequest.Issue.Fields.Status.Name == "Aguardando SD" && (customFieldName == "Service Desk" || customFieldValue == "Service Desk") {
		if jiraRequest.Issue.Fields.Assignee.User == "" || validateSLAUser(jiraRequest.Issue.Fields.Assignee.User) {
			return true
		}
	}

	return false
}

func GetNocUsers() {
	dataUsers := gcp.GetObject("noc-paygo", "jira-users.json")

	json.Unmarshal(dataUsers, &nocUsers)
}

func ValidateRequest(jiraRequest entity.JiraRequest) bool {
	if jiraRequest.User.Name == "Automation for Jira" || jiraRequest.User.Name == "ScriptRunner for Jira" {
		return false
	}

	if jiraRequest.EventName == "issue_property_set" {
		return false
	}

	if jiraRequest.Issue.Key == "" {
		return false
	}

	return true
}
