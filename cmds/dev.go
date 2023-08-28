// contains private commands, not available to the public
package cmds

import (
	"github.com/bwmarrin/discordgo"
	"github.com/rexagod/newman/utils"
)

func SyncCMD(s *discordgo.Session, i *discordgo.InteractionCreate) {

	isDev := utils.CheckIfDev(i.Member.User.ID)

	if isDev != true {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Only daddy gets to use this!",
			}})

	} else {

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource})

		// utils.RegisterCommands(s, "")
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "Synced!",
		})
	}
}
