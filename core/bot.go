package core

import (
	"log"
	"os"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/conf"
	"github.com/rexagod/newman/utils"
)

var (
	Dg *discordgo.Session
)

func DiscordConnect() {
	log.Println("Starting connection to Discord.")
	var err error
	Dg, err = discordgo.New(os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatal("Sed, Couldn't connect to Discord,", err)
		return
	}

	Dg.AddHandler(InteractionHandler)

	Dg.Identify.Intents |= discordgo.IntentsAll

	Dg.StateEnabled = false
	Dg.LogLevel = discordgo.LogError
	Dg.SyncEvents = true

	err = Dg.Open()
	if err != nil {
		log.Fatal("Sed, Couldn't connect to Discord,", err)
	}

	log.Println("Yo! Connected to discord as " + Dg.State.User.Username + "[" + Dg.State.User.ID + "]")

	apps, _ := Dg.ApplicationCommands(Dg.State.Application.ID, "")
	conf.ApplicationCommands = apps

	env := os.Getenv("ENV_TYPE")
	log.Println("Loaded ENV as", env)

	syncTree := os.Getenv("SYNC_TREE")
	toSync, err := strconv.ParseBool(syncTree)

	if err != nil {
		log.Fatal("Error parsing SYNC_TREE env:", err)
	}

	if toSync {
		log.Println("Auto Tree syncing is enabled, syncing now..")
		utils.RegisterCommands(Dg, "")
	} else {
		log.Println("Auto Tree syncing is disabled, skipping!")
	}

}
