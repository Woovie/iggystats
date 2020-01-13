package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const pChatInfoSeparator = "DO NOT MODIFY FROM THIS POINT ONWARDS:"
const pChatInstructions = "- Use `!pchat start` to re-invite the user to this channel." + `
` + "- Use `!pchat stop` to remove the user." + `
` + "- The user has automatically been invited to this channel."

type pchatInfo struct {
	UserID string
}

func (b *bot) privateChatAction(s *discordgo.Session, m *discordgo.Message, parts []string) {
	fmt.Println(strings.Join(parts, ","))
	if parts[0] == "start" || parts[0] == "stop" {
		var info pchatInfo

		channel, err := s.State.Channel(m.ChannelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "ERROR: Could not load channel info: "+err.Error())
			return
		}

		// Check parent
		if channel.ParentID != pchatCategory {
			s.ChannelMessageSend(m.ChannelID, "DENIED: This channel is not inside the p-chat category")
			return
		}

		// Extract info from topic
		{
			list := strings.Split(channel.Topic, pChatInfoSeparator)
			if len(list) != 2 {
				s.ChannelMessageSend(m.ChannelID, "ERROR: Topic is malformed")
				return
			}

			err := json.Unmarshal([]byte(list[1]), &info)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "ERROR: Topic is malformed")
				return
			}
		}

		// Ensure the dude exists
		_, err = b.Member(m.GuildID, info.UserID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "ERROR: could not retrieve target user: "+err.Error())
			return
		}

		if parts[0] == "start" {
			err = s.ChannelPermissionSet(m.ChannelID, info.UserID, "member", discordgo.PermissionReadMessages, 0)
		} else if parts[0] == "stop" {
			err = s.ChannelPermissionDelete(m.ChannelID, info.UserID)
		}

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "ERROR: Could not set target user channel permissions: "+err.Error())
			return
		}

		b.okHand(m)

		return
	}

	targetUser := parts[0]
	if !userRegexp.MatchString(targetUser) {
		return
	}

	targetUID := userRegexp.FindStringSubmatch(targetUser)[1]

	target, err := b.Member(m.GuildID, targetUID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "ERROR: Could not get target user guild member: "+err.Error())
		return
	}

	info := pchatInfo{targetUID}
	infoBytes, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	channel, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
		Name:     target.User.Username + "-" + target.User.Discriminator,
		Type:     discordgo.ChannelTypeGuildText,
		Topic:    pChatInstructions + "\n\n\n" + pChatInfoSeparator + string(infoBytes),
		ParentID: pchatCategory,
		NSFW:     false,
	})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "ERROR: Could not create pchat channel: "+err.Error())
		return
	}

	s.ChannelMessageSend(channel.ID, pChatInstructions)

	err = s.ChannelPermissionSet(channel.ID, info.UserID, "member", discordgo.PermissionReadMessages, 0)
	if err != nil {
		s.ChannelMessageSend(channel.ID, "ERROR: could not give user read permission to channel: "+err.Error())
	}
}
