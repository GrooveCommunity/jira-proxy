module github.com/GrooveCommunity/proxy-jira

go 1.16

require (
	cloud.google.com/go/pubsub v1.11.0 // indirect
	github.com/gorilla/mux v1.8.0
	google.golang.org/api v0.47.0 // indirect
	
	github.com/GrooveCommunity/glib-noc-event-structs v0.0.0
)

replace github.com/GrooveCommunity/glib-noc-event-structs v0.0.0 => /go/src/github.com/GrooveCommunity/glib-noc-event-structs
