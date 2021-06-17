package internal

import (
	"context"

	"log"

	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/GrooveCommunity/proxy-jira/entity"
)

func PublicMessage(projectID, topicName string, jiraRequest entity.JiraRequest) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalln("Ocorreu um erro na comunicação com o PubSub. \nError: " + err.Error())
	}

	defer client.Close()

	topic := client.Topic(topicName)
	defer topic.Stop()

	jiraRequestJson, errJson := json.Marshal(jiraRequest)
	if errJson != nil {
		log.Fatalln("Ocorreu um erro na conversão de json para o jiraRequest. \nError: " + err.Error())
	}

	resp := topic.Publish(ctx, &pubsub.Message{
		Data: jiraRequestJson,
	})

	msgID, errPublish := resp.Get(ctx)
	if errPublish != nil {
		log.Fatalln("Ocorreu um erro na publicação para o topic " + topicName + ". \nError: " + err.Error())
	}

	log.Println("Mensagem " + msgID + " criada com sucesso!")
}
