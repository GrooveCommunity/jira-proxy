package internal

import (
	"context"

	"log"

	"cloud.google.com/go/pubsub"
	"github.com/GrooveCommunity/glib-noc-event-structs/entity"
)

func PublicMessage(projectID, topicName string, payload []byte) {

	var err error

	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	defer client.Close()

	if err != nil {
		log.Fatal(entity.ResponseError{
			Message:    "Ocorreu um erro na comunicação com o PubSub.",
			StatusCode: 504,
			Error:      err,
		})
	}

	topic := client.Topic(topicName)
	defer topic.Stop()

	resp := topic.Publish(ctx, &pubsub.Message{
		Data: payload,
	})

	msgID, errPublish := resp.Get(ctx)
	if errPublish != nil {
		log.Fatal(entity.ResponseError{
			Message:    "Ocorreu um erro na publicação para o topic " + topicName,
			StatusCode: 500,
			Error:      err,
		})
	}

	log.Println("Mensagem " + msgID + " no tópico " + topicName + " criada com sucesso!")
}
