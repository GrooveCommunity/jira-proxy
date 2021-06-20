package internal

import (
	"encoding/json"
	"reflect"

	"log"
	"strings"
	"time"

	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

type customFields map[string]interface{}

func ForwardIssue(jiraRequest entity.JiraRequest, body []byte, projectID, topicDispatcher, topicMetrics string) {

	var customFieldsData customFields

	json.Unmarshal(body, &customFieldsData)

	jiraCustomfields := unmarchallMap(customFieldsData)

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
		CustomFields: jiraCustomfields,
		Transitions:  jiraTransitions,
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
	}

	payload, errPayLoad := json.Marshal(jiraIssue)

	if errPayLoad != nil {
		log.Fatal(entity.ResponseError{
			Message:    "Erro na conversão do payload para JSON",
			StatusCode: 500,
			Error:      errPayLoad,
		})
	}

	go PublicMessage(projectID, topicDispatcher, payload)
	go PublicMessage(projectID, topicMetrics, payload)
}

func unmarchallMap(dataMap map[string]interface{}) []entity.JiraCustomField {
	var customFields []entity.JiraCustomField

	for key, item := range dataMap {
		if (key == "issue" || key == "fields") && reflect.TypeOf(item).Kind() == reflect.Map {
			customFields = unmarchallMap(item.(map[string]interface{}))
		}

		if strings.HasPrefix(key, "customfield") && item != nil && reflect.TypeOf(item).Kind() == reflect.String {
			customFields = append(customFields, entity.JiraCustomField{ID: key, Value: item.(string)})
		}
	}

	return customFields
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
