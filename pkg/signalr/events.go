package signalr

type EventType struct {
	NodeEventType string `json:"nodeEventType"`
}

type StakingInfoEvent struct {
	EventType
	Enabled     bool   `json:"enabled"`
	Staking     bool   `json:"staking"`
	Error       string `json:"errors"`
	BlockSize   int    `json:"currentBlockSize"`
	BlockTx     int    `json:"currentBlockTx"`
	Weight      int    `json:"weight"`
	StakeWeight int    `json:"netStakeWeight"`
	Time        int    `json:"expectedTime"`
}

type WalletGeneralInfoEvent struct {
	EventType
	Balances            []Account
	Name                string `json:"walletName"`
	Network             string `json:"network"`
	LastBlockSyncHeight int    `json:"lastBlockSyncedHeight"`
	ChainTip            int    `json:"chainTip"`
	IsChainSynced       bool   `json:"isChainSynced"`
	Nodes               int    `json:"connectedNodes"`
}

type Account struct {
	Name        string `json:"accountName"`
	Confirmed   int    `json:"amountConfirmed"`
	Unconfirmed int    `json:"amountUnconfirmed"`
	Amount      int    `json:"spendableAmount"`
}
