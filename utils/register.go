package utils

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/conf"
)

// TODO: remove all the synced commands before syncing again.

func RegisterCommands(s *discordgo.Session, GuildId string) {
	commands := conf.Commands
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildId, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}

		log.Println("Added Command: " + cmd.Name + " [" + strconv.Itoa(i) + "]")
	}
}
