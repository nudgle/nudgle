package http

import (
	"encoding/json"
	"github.com/nudgle/nudgle/pkg/signalr"
	"net/http"
)

type Controller struct {
	Receiver chan interface{}
}

// Staking is the controller for the /stakinginfo route
func (c Controller) Staking(w http.ResponseWriter, req *http.Request) {
	var p signalr.StakingInfoEvent
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Receiver <- p
}

// Wallet is the controller for the /walletgeneralinfo route
func (c Controller) Wallet(w http.ResponseWriter, req *http.Request) {
	var p signalr.WalletGeneralInfoEvent
	err := json.NewDecoder(req.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.Receiver <- p
}
