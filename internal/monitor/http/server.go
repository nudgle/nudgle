package http

import (
	monitor "github.com/nudgle/nudgle/pkg/monitor/config"
	"net/http"
)

// Start runs an HTTP Server using net/http as a framework
func Start(config *monitor.MonitorServiceConfiguration, receiver chan interface{}) {
	handler := Controller{Receiver: receiver}
	http.HandleFunc("/stakinginfo", handler.Staking)
	http.HandleFunc("/walletgeneralinfo", handler.Wallet)

	http.ListenAndServe(config.Server.Address, nil)
}
