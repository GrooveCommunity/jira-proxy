package internal

import (
	"context"

	"log"

	//	"encoding/json"

	//	"sync"
	//	"sync/atomic"

	"cloud.google.com/go/pubsub"
	"github.com/GrooveCommunity/proxy-jira/entity"
	//	"google.golang.org/api/iterator"
)

func ForwardDispatcher(projectID string, jiraRequest entity.JiraRequest) {
	publicMessage(projectID, "projects/monitoria-groovetech/topics/dispatcher", jiraRequest)
}

func publicMessage(projectID, topicName string, jiraRequest entity.JiraRequest) {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalln("Ocorreu um erro na comunicação com o PubSub. \nError: " + err.Error())
	}

	defer client.Close()

	topic := client.Topic(topicName)
	defer topic.Stop()

	var results []*pubsub.PublishResult
	result := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
	})

	results = append(results, result)

	for _, r := range results {
		id, errResult := r.Get(ctx)
		if errResult != nil {
			log.Fatal(errResult)
		}

		log.Println("Published a message with a message ID: " + id)
	}

	/*topics := client.Topics(ctx)

	for {
		item, errTopics := topics.Next()
		if err == iterator.Done {
			break
		}

		if errTopics != nil {
			log.Fatal(errTopics)
		}

		log.Println(item)

	}

	//log.Println(topics.Next())*/

	/*topic := client.Topic("projects/monitoria-groovetech/topics/dispatcher")

	jiraRequestJson, errJson := json.Marshal(jiraRequest)
	log.Println(jiraRequestJson)

	if errJson != nil {
		log.Fatalln("Ocorreu um erro na conversão de json para o jiraRequest. \nError: " + err.Error())
	}

	resp := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("teste go"),
	})

	msgID, errPublish := resp.Get(ctx)
	if errPublish != nil {
		log.Fatalln("Ocorreu um erro na publicação para o topic " + topicName + ". \nError: " + err.Error())
	}

	log.Println("Mensagem " + msgID + " criada com sucesso!")*/

	//topic.Publish(ctx, msg{})

}
