package cmds

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func PingCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!!",
		},
	})

	if err != nil {
		log.Println("Error responding to ping: ", err)
	}
}
