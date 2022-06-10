package core

import (
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/handler"
	"github.com/rexagod/newman/core/queries"
	"k8s.io/klog/v2"
	"time"
)

func undoDelete() *state.State {
	s := state.New("Bot " + R.token)
	s.PreHandler = handler.New()
	s.PreHandler.AddSyncHandler(func(e *gateway.MessageDeleteEvent) {
		message, err := s.Message(e.ChannelID, e.ID)
		if err != nil {
			klog.Errorf("Failed to get message: %v", err)
		} else {
			who := message.Author.Username
			when := message.Timestamp.Time().Add(time.Hour*5 + time.Minute*30).Format(time.RFC3339)
			what := message.Content
			s, err := addRow(queries.Q[queries.ADDDELETEDMESSAGE], who, when, what)
			if err != nil {
				klog.Errorf("Failed to add row: %v", err)
			}
			klog.Info(s)
		}
	})
	s.AddIntents(gateway.IntentGuildMessages)
	s.AddIntents(gateway.IntentDirectMessages)

	return s
}
