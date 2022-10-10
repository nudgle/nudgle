package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	config2 "nudgle/internal/config"
	"nudgle/internal/signalr"
	strax "nudgle/pkg/strax/client"
)

type MessagePrinterHandler struct {
	onStart func()
}

// Default handles incoming events from SignalR
func (m MessagePrinterHandler) Default(ctx context.Context, target string, args []json.RawMessage) error {
	for _, val := range args {
		event, err := Parse(val)
		if err != nil {
			log.Println(err)
			return err
		}
		Handle(event, val)
	}
	return nil
}

// Start starts listening on the signalR server
// and handles the event messages coming from it
func Start() {
	conf := config2.GetIndexerConfig()
	apiClient := strax.GetClient(conf.Node)
	signalConnection, err := apiClient.GetSignalR()
	if err != nil {
		panic(err)
	}
	connStr := fmt.Sprintf(
		"%s/%s",
		signalConnection.URI,
		conf.Node.Hub,
	)
	hubName := ""
	client, err := signalr.NewClient(connStr, hubName)
	if err != nil {
		log.Println(err)
		return
	}

	listenCtx, cancel := context.WithCancel(context.Background())

	started := make(chan struct{}, 1)
	mph := &MessagePrinterHandler{
		onStart: func() {
			started <- struct{}{}
			cancel()
		},
	}

	go func() {
		err = client.Listen(listenCtx, mph)
		if err != nil {
			log.Println(err)
			started <- struct{}{}
		}
	}()
	//
	log.Println(<-started)
}
