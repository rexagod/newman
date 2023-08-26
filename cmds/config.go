package cmds

import "github.com/bwmarrin/discordgo"

var ApplicationCommands []*discordgo.ApplicationCommand
var RegisterCommand string

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "PONG!",
	},
}
