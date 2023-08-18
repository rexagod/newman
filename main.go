package main

import (
	"os"
	"os/signal"

	"github.com/diamondburned/arikawa/v3/state"
	"k8s.io/klog/v2"

	"github.com/rexagod/newman/core"
	"github.com/rexagod/newman/internal"
)

func main() {
	var l internal.Loader
	err := l.Load()
	if err != nil {
		klog.Fatalf("failed to load metadata: %v", err)
	}

	// Start the bot.
	var s *state.State
	s, err = core.Start(&l)
	if err != nil {
		klog.Fatalf("failed to start bot: %v", err)
	}
	defer func(s *state.State) {
		err := s.Close()
		if err != nil {
			klog.Errorf("failed to close bot: %v", err)
		}
	}(s)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	<-sigs
	klog.Info("received interrupt, shutting down")
}
