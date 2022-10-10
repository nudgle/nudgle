package big_transactions

import (
	"encoding/hex"
	"fmt"
	"github.com/nudgle/nudgle/pkg/monitor/config"
	"github.com/nudgle/nudgle/pkg/monitor/filter"
	strax "github.com/nudgle/nudgle/pkg/strax/client"
	"github.com/nudgle/nudgle/pkg/strax/model"
	"log"
	"strconv"
	"strings"
	"time"
)

func init() {
	hwPlugin := &BigTransactionFilter{}
	filter.Register(hwPlugin.GetCmdName(), hwPlugin)
}

type BigTransactionFilter struct {
	filter.Handler
}

func (e *BigTransactionFilter) GetCmdName() string {
	return "detect-big-transactions"
}

// Execute gets called by the worker in a chain with a transaction
// id argument. It then uses the strax api to fetch transaction data
// from a full node
func (e *BigTransactionFilter) Execute(
	config *config.MonitorServiceConfiguration,
	tx string,
) string {
	client := strax.GetClient(config.Node)
	transaction, err := client.GetTx(tx)
	if err != nil {
		log.Println(err)
		return ""
	}
	output := ""
	output += fmt.Sprintf("Transaction occurred: \n")
	output += fmt.Sprintf("TXID: %s\n", transaction.ID)

	// fetching input addresses
	vinAddresses := e.getInputAddresses(client, transaction.Vin)

	if len(vinAddresses) > 0 && vinAddresses[0] == "yU2jNwiac7XF8rQvSk2bgibmwsNLkkhsHV" {
		log.Println("Skipping withdraw transaction")
		return ""
	}

	date, _ := strconv.ParseInt(strconv.Itoa(transaction.Time), 10, 64)
	output += fmt.Sprintf("Time: %s\n", time.Unix(date, 0))
	scriptHash, err := strax.GetVout(transaction.Vout, "scripthash")
	if err != nil {
		log.Println(err)
		return ""
	}
	if scriptHash.ScriptPubKey.Addresses[0] != "yU2jNwiac7XF8rQvSk2bgibmwsNLkkhsHV" {
		log.Println("Skipping non cross-chain transaction")
		return ""
	}
	returnAddress, err := e.getReturnAddress(transaction.Vout)
	if err != nil {
		log.Println(err)
		return ""
	}

	output += fmt.Sprintf("Amount: %.2f\n", scriptHash.Value)
	output += fmt.Sprintf("Return address: %s\n", returnAddress)
	return output
}

func (e *BigTransactionFilter) getReturnAddress(vout []model.Vout) ([]byte, error) {
	opReturn, err := strax.GetVout(vout, "nulldata")
	if err != nil {
		log.Println(err)
	}
	asm := opReturn.ScriptPubKey.Asm
	arguments := strings.Split(asm, " ")
	if arguments[0] != "OP_RETURN" {
		log.Println("Not an OP_RETURN vout")
		return nil, fmt.Errorf("Not an OP_RETURN vout")
	}
	returnAddress, err := hex.DecodeString(arguments[1])
	if err != nil {
		log.Println("Failed to decode hex")
		return nil, fmt.Errorf("Failed to decode hex")
	}
	return returnAddress, nil
}

func (e *BigTransactionFilter) getInputAddresses(
	client *strax.Client,
	inputs []model.Vin,
) []string {
	addresses := []string{}
	for _, value := range inputs {
		if value.ID == "" {
			continue
		}
		tx, err := client.GetTx(value.ID)
		if err != nil {
			log.Println("Failed to fetch tx", err)
			return []string{}
		}
		addresses = append(addresses, tx.Vout[value.Vout].ScriptPubKey.Addresses[0])
	}
	return addresses
}
