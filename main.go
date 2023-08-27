package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rexagod/newman/core"
	"github.com/rexagod/newman/database"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Lol! I couldn't load the fkin env")
	} else {
		log.Println("Env loaded, my mane! Good work.")
	}

	err = database.Initialize()
	if err != nil {
		log.Panic(err)
	}

	err = database.CreateTables()
	if err != nil {
		log.Fatal(err)
	}

	core.DiscordConnect()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
