package cmds

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/database"
)

func PingCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!!",
		},
	})

	if err != nil {
		log.Println("Error responding to ping: ", err)
	}
}

type SnipeData struct {
	MessageID int64
	ChannelID int64
	AuthorID  int64
	Content   json.RawMessage
	DeletedAt time.Time
}

func SnipeCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {

	toRetrieve := 1
	if len(i.ApplicationCommandData().Options) > 0 {
		toRetrieve = int(i.ApplicationCommandData().Options[0].IntValue())
	}

	query := `SELECT * FROM snipes WHERE channel_id = $1 ORDER BY deleted_at DESC LIMIT $2`

	rows, err := database.GetDB().Query(query, i.ChannelID, toRetrieve)

	if err != nil {
		log.Println("Error querying the snipes:", err)
		return
	}

	defer rows.Close()

	var snipes []SnipeData
	for rows.Next() {
		var snipe SnipeData
		err := rows.Scan(&snipe.MessageID, &snipe.ChannelID, &snipe.AuthorID, &snipe.Content, &snipe.DeletedAt)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}
		snipes = append(snipes, snipe)
	}

	if len(snipes) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Nothing to snipy snipe :c",
			}})
	}

	var embeds []*discordgo.MessageEmbed
	for _, snipe := range snipes {
		var contentData map[string]interface{}
		err := json.Unmarshal(snipe.Content, &contentData)

		if err != nil {
			log.Println("Error unmarshaling content JSON:", err)
			continue
		}

		messageContent := contentData["text"].(string)
		attachmentURL, attachmentExists := contentData["attachment"].(string)
		if !attachmentExists {
		}

		embed := &discordgo.MessageEmbed{
			Description: fmt.Sprintf(
				"Author: <@%d>\n**__Message Content__**\n%s", snipe.AuthorID, messageContent,
			),
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Deleted at ",
			},
			Timestamp: snipe.DeletedAt.Format(time.RFC3339),
			Image: &discordgo.MessageEmbedImage{
				URL: attachmentURL,
			},
			Color: 0x2F3136,
		}
		embeds = append(embeds, embed)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		}})

}
