package indexer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nudgle/nudgle/internal/cache"
	config2 "github.com/nudgle/nudgle/internal/config"
	"github.com/nudgle/nudgle/pkg/signalr"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Event string

// Parse takes in a raw json argument and decodes it
// The function detects which Event type the event is
// and returns it
func Parse(event json.RawMessage) (Event, error) {
	var data signalr.EventType
	json.Unmarshal(event, &data)
	opts := strings.Split(data.NodeEventType, ",")
	if len(opts[0]) < 1 {
		return Event(""), errors.New("Event type not set")
	}
	switch opts[0] {
	case "Stratis.Bitcoin.Features.SignalR.Events.StakingInfo":
		return Event("StakingInfo"), nil
		break
	case "Stratis.Bitcoin.Features.SignalR.Events.WalletGeneralInfo":
		return Event("WalletGeneralInfo"), nil
		break
	default:
		return Event(""), errors.New("Unknown event type")
		break
	}
	return Event(""), errors.New("Unknown event type")
}

// ParseStakingInfo json decodes the raw json message
// into a signalr.StakingInfoEvent struct
func ParseStakingInfo(data json.RawMessage) *signalr.StakingInfoEvent {
	var payload signalr.StakingInfoEvent
	json.Unmarshal(data, &payload)
	return &payload
}

// ParseWalletGeneralInfo json decodes the raw json message
// into a signalr.WalletGeneralInfoEvent struct
func ParseWalletGeneralInfo(data json.RawMessage) *signalr.WalletGeneralInfoEvent {
	var payload signalr.WalletGeneralInfoEvent
	json.Unmarshal(data, &payload)
	return &payload
}

// Handle starts processing events based on their type
func Handle(event Event, data json.RawMessage) {
	switch event {
	case Event("StakingInfo"):
		data := ParseStakingInfo(data)
		ProcessStakingInfo(data)
		break
	case Event("WalletGeneralInfo"):
		data := ParseWalletGeneralInfo(data)
		ProcessWalletGeneralInfo(data)
		break
	}
}

// SendData dispatches the event to the monitor service
// for further analysis of the block
func SendData(path string, data interface{}) {
	config := config2.GetIndexerConfig()
	endpoint := fmt.Sprintf(
		"http://%s:%d/%s",
		config.Monitor.Host,
		config.Monitor.Port,
		path,
	)
	log.Println("HTTP JSON POST URL:", endpoint)

	var jsonData, _ = json.Marshal(data)
	request, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{}
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		log.Println(err)
		return
	}

	_, _ = ioutil.ReadAll(response.Body)
}

func ProcessStakingInfo(data *signalr.StakingInfoEvent) {
	// TODO: Implement event handling
}

// ProcessWalletGeneralInfo takes in the event struct, it then
// checks if it has been processed before, to avoid sending the
// same block twice to the monitor app
func ProcessWalletGeneralInfo(data *signalr.WalletGeneralInfoEvent) {
	memory := cache.Memory()
	result, ok := memory.Get("lastBlockSyncedHeight")
	if !ok {
		log.Println(data.LastBlockSyncHeight)
		memory.Set("lastBlockSyncedHeight", data.LastBlockSyncHeight, cache.NoExpiration)
		SendData("walletgeneralinfo", data)
		return
	}
	if result.(int) < data.LastBlockSyncHeight {
		memory.Set("lastBlockSyncedHeight", data.LastBlockSyncHeight, cache.NoExpiration)
		SendData("walletgeneralinfo", data)
	}
}
