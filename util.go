package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func isModRole(role string) bool {
	for _, s := range modRoles {
		if s == role {
			return true
		}
	}
	return false
}

func composeMessageURL(m *discordgo.Message) string {
	return "https://discordapp.com/channels/" + m.GuildID + "/" + m.ChannelID + "/" + m.ID
}

func channelsRemove(slice []*discordgo.Channel, cid string) []*discordgo.Channel {
	var i int
	for index, v := range slice {
		if v.ID == cid {
			i = index
			break
		}
	}

	return append(slice[:i], slice[i+1:]...)
}

func channelsInsertAfter(s []*discordgo.Channel, after string, value *discordgo.Channel) []*discordgo.Channel {
	var i int
	for index, v := range s {
		if v.ID == after {
			i = index + 1
			break
		}
	}

	if len(s) == i {
		return append(s, value)
	}

	s = append(s, nil)
	copy(s[i+1:], s[i:])
	s[i] = value
	return s
}

func printChannels(chans []*discordgo.Channel) {
	fmt.Println("channels: {")
	for _, c := range chans {
		fmt.Printf("\t%s @ %d (%s)\n", c.Name, c.Position, c.ID)
	}
	fmt.Println("}")
}

// MemberName is a standardised way to get the name of a member
func MemberName(member *discordgo.Member) string {
	if member.Nick != "" {
		return member.Nick
	}
	return member.User.Username
}

func everyone(guildID string) string {
	return "<@&" + guildID + ">"
}

func stripEveryone(guildID string, message string) string {
	message = strings.ReplaceAll(message, everyone(guildID), "")
	message = strings.ReplaceAll(message, "@everyone", "")
	message = strings.ReplaceAll(message, "@here", "")
	return message
}
