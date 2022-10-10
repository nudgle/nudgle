package main

import (
	"github.com/nudgle/nudgle/internal/config"
	"github.com/nudgle/nudgle/internal/monitor/discord"
	"github.com/nudgle/nudgle/internal/monitor/http"
	"github.com/nudgle/nudgle/internal/monitor/processor"
	"os"
)

func main() {
	config.GetFlags(config.MonitorService)
	config := config.GetMonitorConfig()
	started := make(chan os.Signal, 1)
	bot := discord.NewBot(config)
	go bot.Listen()
	workers := processor.New(config, bot.Channel)
	workers.Start()
	go http.Start(config, workers.Receiver)
	<-started
}
