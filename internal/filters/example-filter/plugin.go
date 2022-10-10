package example_filter

import (
	"log"
	"nudgle/pkg/monitor/config"
	"nudgle/pkg/monitor/filter"
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
