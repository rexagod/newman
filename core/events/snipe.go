package events

import (
	"encoding/json"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/database"
)

// Stores deleted messages to later on use with snipe command.
func SnipeEvent(s *discordgo.Session, m *discordgo.MessageDelete) {

	if m.BeforeDelete == nil {
		return // means the bot didn't have the message cached in the state.
	}

	if m.BeforeDelete.GuildID == "" || m.BeforeDelete.Author.Bot {
		return
	}

	attachmentURL := ""
	if len(m.BeforeDelete.Attachments) > 0 {
		if m.BeforeDelete.Attachments[0].ContentType == "image/png" || m.BeforeDelete.Attachments[0].ContentType == "image/jpeg" {
			attachmentURL = m.BeforeDelete.Attachments[0].ProxyURL
		}

	}
	contentMap := map[string]interface{}{
		"text":       m.BeforeDelete.Content,
		"attachment": attachmentURL,
	}
	contentJSON, err := json.Marshal(contentMap)
	if err != nil {
		log.Println("[Snipe] Error marshaling content map to JSON: ", err)
		return
	}

	query := `
	INSERT INTO snipes (message_id, channel_id, author_id, content)
	VALUES ($1, $2, $3, $4)
	`
	database.GetDB().Exec(query, m.ID, m.ChannelID, m.BeforeDelete.ID, contentJSON)

}
