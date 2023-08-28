package conf

import "github.com/bwmarrin/discordgo"

var ApplicationCommands []*discordgo.ApplicationCommand
var RegisterCommand string

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "PONG!",
	},
	{
		Name:        "sync",
		Description: "Syncs commands with discord",
	},
}
