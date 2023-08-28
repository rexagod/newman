package core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/cmds"
)

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}

	if i.Data.Type() != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		cmds.PingCMD(s, i)
	case "sync":
		cmds.SyncCMD(s, i)
	}

}
