package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/rexagod/newman/core"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Lol! I couldn't load the fkin env")
	} else {
		log.Println("Env loaded, my mane! Good work.")
	}

	core.DiscordConnect()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

}
