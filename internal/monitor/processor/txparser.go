package processor

import (
	strax "github.com/nudgle/nudgle/pkg/strax/client"
	"github.com/nudgle/nudgle/pkg/strax/model"
	"log"
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
