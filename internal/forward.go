package internal

import (
	"log"
	//	glibEntity "github.com/GrooveCommunity/glib-noc-event-structs/entity"
	"encoding/json"
	"reflect"

	"strings"

	"github.com/GrooveCommunity/proxy-jira/entity"

	glibentity "github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

type customFields map[string]interface{}

func ForwardIssue(jiraRequest entity.JiraRequest, body []byte) {

	var customFieldsData customFields

	json.Unmarshal(body, &customFieldsData)

	var jiraCustomfields []glibentity.JiraCustomField

	jiraCustomfields = unmarchallMap(customFieldsData)

	log.Println(jiraCustomfields)
}

func ForwardDispatcher(projectID string, jiraRequest entity.JiraRequest) {
	PublicMessage(projectID, "dispatcher-jira-paygo", jiraRequest)
}

func unmarchallMap(dataMap map[string]interface{}) []glibentity.JiraCustomField {
	var customFields []glibentity.JiraCustomField

	for key, item := range dataMap {
		if (key == "issue" || key == "fields") && reflect.TypeOf(item).Kind() == reflect.Map {
			customFields = unmarchallMap(item.(map[string]interface{}))
		}

		if strings.HasPrefix(key, "customfield") && item != nil && reflect.TypeOf(item).Kind() == reflect.String {
			customFields = append(customFields, glibentity.JiraCustomField{ID: key, Value: item.(string)})
		}
	}

	return customFields
}
