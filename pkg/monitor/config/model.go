package config

import "github.com/nudgle/nudgle/pkg/strax/model"

type (
	MonitorServiceConfiguration struct {
		Bot    Bot        `fig:"bot" validate:"required"`
		Node   model.Node `fig:"node" validate:"required"`
		Server Server     `fig:"server" validate:"required"`
	}

	Server struct {
		Address       string `fig:"address" default:":8090"`
		WorkerThreads int    `fig:"workerThreads" default:"2"`
	}

	Bot struct {
		Token     string `fig:"token" validate:"required"`
		ChannelID string `fig:"channelId" validate:"required"`
	}
)
