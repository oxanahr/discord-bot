package utils

import (
	"github.com/getsentry/sentry-go"
	"github.com/oxanahr/discord-bot/cmd/context"
	"log"
)

// SendChannelMessage sends a channel message to channel with channel id equal to m.ChannelID
func SendChannelMessage(channelID string, message string) {
	_, err := context.Dg.ChannelMessageSend(channelID, message)
	if err != nil {
		log.Fatalln(err)
	}
}

func SendPrivateMessage(userID string, message string) {
	channel, err := context.Dg.UserChannelCreate(userID)
	if err != nil {
		log.Fatalln(err)
		return
	}
	SendChannelMessage(channel.ID, message)
}

// Mention Returns the passed username with a prefix of @
func Mention(userID string) string {
	user, err := context.Dg.User(userID)
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}
	return user.Mention()
}

// Username	Returns only the username
func Username(userID string) string {
	user, err := context.Dg.User(userID)
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}
	return user.Username
}

// Padding Calculates the offset needed to reach the maximum column width
func Padding(s string, p int) int {
	padding := 0
	if len(s)%p != 0 {
		padding = p * (len(s)/20 + 1)
	}
	return padding
}
