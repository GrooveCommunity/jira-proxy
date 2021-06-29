package internal

import (
	"log"
	"os"

	"time"

	"github.com/bwmarrin/discordgo"
)

func SendMessageToChannel(url, issueKey, message string, color int) {
	token := os.Getenv("TOKEN_DISPATCHER_PAYGO_DISCORD")
	channelID := os.Getenv("CHANNEL_ID_DISCORD")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("erro ao abrir conexão com o servidor")
		panic(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	msg := discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeRich,
		Title:       issueKey,
		URL:         url,
		Description: message,
		Timestamp:   time.Now().Format(time.RFC3339),
		Color:       color,
	}

	_, errMsg := dg.ChannelMessageSendEmbed(channelID, &msg)

	if errMsg != nil {
		panic(errMsg)
	}
}
