package processor

import (
	"log"
	strax "nudgle/pkg/strax/client"
	"nudgle/pkg/strax/model"
)

// ProcessWalletGeneralInfo takes in the latest block height
// It then fetches the block data and transactions to parse
// Then analyses each transaction separately
func (p *Processor) ProcessWalletGeneralInfo(lastBlockSyncHeight int) {
	log.Println("LastHeight", lastBlockSyncHeight)
	client := strax.GetClient(p.Config.Node)
	hash, err := client.GetBlockHash(lastBlockSyncHeight)
	if err != nil {
		log.Println("Failed to fetch block hash", err)
		return
	}
	block, err := client.GetBlock(hash)
	if err != nil {
		log.Println("Failed to fetch block data", err)
		return
	}
	for _, tx := range block.Transactions {
		txList := model.TransactionList{Transaction: tx}
		p.Receiver <- txList
	}
}
