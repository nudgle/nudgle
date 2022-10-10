// Package discord listens on a go channel for messages to be forwarded
// to a discord channel
package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	monitor "nudgle/pkg/monitor/config"
)

type Discord struct {
	Channel   chan string
	ChannelID string
	Session   *discordgo.Session
}

// NewBot returns a reference of a new Discord struct
// which contains a bot token, a channel to communicate
// the messages to be sent to the channel, and the channelId
func NewBot(config *monitor.MonitorServiceConfiguration) *Discord {
	token := fmt.Sprintf("Bot %s", config.Bot.Token)
	session, err := discordgo.New(token)
	if err != nil {
		log.Println(err)
		return nil
	}
	return &Discord{
		Channel:   make(chan string),
		ChannelID: config.Bot.ChannelID,
		Session:   session,
	}
}

// Listen starts listening on the message channel
// It sends a message to the discord channel for every
// event passed to the go chan
func (d *Discord) Listen() {
	for message := range d.Channel {
		d.SendMessage(message)
	}
}

// SendMessage forwards the text message to the discord channel
// that has been configured in the config.yaml
func (d *Discord) SendMessage(text string) {
	_, err := d.Session.ChannelMessageSend(d.ChannelID, text)
	if err != nil {
		log.Println(err)
		return
	}
}
