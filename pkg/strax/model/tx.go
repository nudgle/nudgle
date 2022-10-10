package model

type (
	Node struct {
		Scheme      string `fig:"scheme" default:"http"`
		Host        string `fig:"host" validate:"required"`
		Port        int    `fig:"port" default:"17103"`
		SignalrPort int    `fig:"signalrPort" default:"17102"`
		Hub         string `fig:"hub" default:"events-hub"`
	}
	Block struct {
		Hash          string   `json:"hash"`
		Size          int      `json:"size"`
		Weight        int      `json:"weight"`
		Height        int      `json:"height"`
		Confirmations int      `json:"confirmations"`
		Time          int      `json:"time"`
		Transactions  []string `json:"tx"`
	}
	Transaction struct {
		Hex           string `json:"hex"`
		ID            string `json:"txid"`
		Hash          string `json:"hash"`
		Size          int    `json:"size"`
		Weight        int    `json:"weight"`
		Vin           []Vin  `json:"vin"`
		Vout          []Vout `json:"vout"`
		Confirmations int    `json:"confirmations"`
		Time          int    `json:"time"`
	}

	TransactionList struct {
		Transaction string
	}

	Vin struct {
		ID       string `json:"txid"`
		Vout     int    `json:"vout"`
		Sequence int    `json:"sequence"`
	}

	Vout struct {
		Value        float64 `json:"value"`
		N            int     `json:"n"`
		ScriptPubKey ScriptPubKey
	}

	ScriptPubKey struct {
		Asm       string   `json:"asm"`
		Hex       string   `json:"hex"`
		Type      string   `json:"type"`
		Addresses []string `json:"addresses,omitempty"`
	}

	SignalR struct {
		URI  string `json:"signalRUri"`
		Port int    `json:"signalRPort"`
	}
)
