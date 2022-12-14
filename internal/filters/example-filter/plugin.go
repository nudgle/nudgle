package example_filter

import (
	"github.com/nudgle/nudgle/pkg/monitor/config"
	"github.com/nudgle/nudgle/pkg/monitor/filter"
	"log"
)

func init() {
	hwPlugin := &ExampleFilter{}
	filter.Register(hwPlugin.GetCmdName(), hwPlugin)
}

type ExampleFilter struct {
	filter.Handler
}

func (e *ExampleFilter) GetCmdName() string {
	return "example-filter"
}

func (e *ExampleFilter) Execute(config *config.MonitorServiceConfiguration, tx string) string {
	log.Println(config, tx)
	return ""
}
