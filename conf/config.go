package conf

import "github.com/bwmarrin/discordgo"

var ApplicationCommands []*discordgo.ApplicationCommand
var RegisterCommand string
var MinValueSnipe = 1.0

var Commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "PONG!",
	},
	{
		Name:        "sync",
		Description: "Syncs commands with discord",
	},
	{
		Name:        "snipe",
		Description: "Retrieves deleted messages.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "count",
				Description: "How many snipes to retrieve? (1 - 10)",
				Required:    false,
				MinValue:    &MinValueSnipe,
				MaxValue:    10.0,
			},
		},
	},
}
