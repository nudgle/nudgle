// Package processor provides the necessary means to run
// multiple transaction analysis in parallel
package processor

import (
	"log"
	monitor "nudgle/pkg/monitor/config"
	"nudgle/pkg/monitor/filter"
	"nudgle/pkg/signalr"
	"nudgle/pkg/strax/model"
	"reflect"
)

type Processor struct {
	Threads  int
	Receiver chan interface{}
	Sender   chan string
	Config   *monitor.MonitorServiceConfiguration
}

// New returns a Processor struct that contains how many
// workers should be running, a job receiver channel
// and the sender channel which sends messages to the
// discord worker
func New(config *monitor.MonitorServiceConfiguration, botChannel chan string) *Processor {
	return &Processor{
		Threads:  config.Server.WorkerThreads,
		Receiver: make(chan interface{}),
		Sender:   botChannel,
		Config:   config,
	}
}

// Start runs the X number of workers as goroutines
func (p *Processor) Start() {
	for i := 1; i <= p.Threads; i++ {
		go p.worker(i)
	}
}

// worker is the individual worker that listens on
// the receiver channel for jobs
func (p *Processor) worker(id int) {
	pluginManager := filter.GetPluginManager()
	log.Println("Worker: ", id)
	for j := range p.Receiver {
		ion := reflect.TypeOf(j)
		log.Println("Job: ", ion.Name(), " taken by worker: ", id)
		switch ion.Name() {
		case "StakingInfoEvent":
			//payload := j.(signalr.StakingInfoEvent)
			break
		case "WalletGeneralInfoEvent":
			payload := j.(signalr.WalletGeneralInfoEvent)
			log.Println(payload.LastBlockSyncHeight)
			p.ProcessWalletGeneralInfo(payload.LastBlockSyncHeight)
			break
		case "TransactionList":
			payload := j.(model.TransactionList)
			pluginManager.Chain(p.Config, &p.Sender, payload.Transaction)
			break
		}

		log.Println("Job: ", ion.Name(), " done by worker: ", id)
	}
}
