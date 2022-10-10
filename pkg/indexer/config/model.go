package config

import "github.com/nudgle/nudgle/pkg/strax/model"

type (
	IndexerServiceConfiguration struct {
		Node    model.Node `fig:"node" validate:"required"`
		Monitor Monitor    `fig:"monitor"`
	}

	Monitor struct {
		Host string `fig:"host" default:"localhost"`
		Port int    `fig:"port" default:"8090"`
	}
)
