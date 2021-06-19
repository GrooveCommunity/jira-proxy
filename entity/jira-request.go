package entity

/*type JiraEvent struct {
	DateEvent time.Time
}*/

type JiraUser struct {
	Name string `json:"displayName"`
}

type JiraIssue struct {
	ID     string
	Key    string
	Fields JiraFields
}

type ChangeLogEvent struct {
	Field string
	From  string `json:"fromString"`
	To    string `json:"toString"`
}

type JiraChangeLog struct {
	Changes []ChangeLogEvent `json:"items"`
}

type JiraPriority struct {
	Name string
}

type JiraResolution struct {
	Name        string
	Description string
}

type JiraAssignee struct {
	User string `json:"displayName"`
}

type JiraStatus struct {
	Name string
}

type JiraCreator struct {
	Name string `json:"displayName"`
}

type JiraReporter struct {
	Name string `json:"displayName"`
}

type JiraType struct {
	Name string
}

type JiraProject struct {
	Name string
}

type JiraFields struct {
	ChangeDate  string `json:"statuscategorychangedate"`
	User        string `json:"displayName"`
	Priority    JiraPriority
	Resolution  JiraResolution
	Assignee    JiraAssignee
	Creator     JiraCreator
	Reporter    JiraReporter
	IssueType   JiraType
	Project     JiraProject
	Created     string
	Updated     string
	Summary     string
	Status      JiraStatus
	Description string
}

type JiraRequest struct {
	EventName string `json:"webhookEvent"`
	User      JiraUser
	Issue     JiraIssue
	ChangeLog JiraChangeLog
}
