package internal

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToChannel(message string) {
	token := os.Getenv("TOKEN_DISPATCHER_PAYGO_DISCORD")
	channelID := os.Getenv("CHANNEL_ID_DISCORD")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("erro ao abrir conexão com o servidor")
		panic(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	_, errMsg := dg.ChannelMessageSend(channelID, message)

	if errMsg != nil {
		panic(errMsg)
	}
}
