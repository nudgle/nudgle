package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"nudgle/pkg/strax/model"
	"strconv"
)

type Client struct {
	Endpoint string
	api      *resty.Client
}

func getResty() *resty.Client {
	// Create a Resty Client
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}

// GetClient takes in a model.Node instance that
// contains information about the endpoint
// and returns an instance of Client
func GetClient(node model.Node) *Client {
	baseUri := fmt.Sprintf(
		"%s://%s:%d",
		node.Scheme,
		node.Host,
		node.Port,
	)
	return &Client{
		Endpoint: baseUri,
		api:      getResty(),
	}
}

func (c *Client) GetSignalR() (*model.SignalR, error) {
	// api/SignalR/getConnectionInfo
	endpoint := fmt.Sprintf("%s/api/SignalR/getConnectionInfo", c.Endpoint)
	resp, err := c.api.R().
		SetHeader("Accept", "application/json").
		Get(endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var signalr model.SignalR
	err = json.Unmarshal(resp.Body(), &signalr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &signalr, nil
}

// GetBlockTransactions returns all the transactions within
// a block
func (c *Client) GetBlockTransactions(height int) []string {
	hash, err := c.GetBlockHash(height)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	block, err := c.GetBlock(hash)
	if err != nil {
		log.Println(err)
		return []string{}
	}
	return block.Transactions
}

// GetBlockHash takes in an integer and returns the hash
// of that block as a string
func (c *Client) GetBlockHash(height int) (string, error) {
	endpoint := fmt.Sprintf("%s/api/Consensus/getblockhash", c.Endpoint)
	resp, err := c.api.R().
		SetQueryParams(map[string]string{
			"height": strconv.Itoa(height),
		}).
		SetHeader("Accept", "application/json").
		Get(endpoint)
	if err != nil {
		log.Println(err)
		return "", err
	}
	hash, err := strconv.Unquote(string(resp.Body()))
	if err != nil {
		log.Println(err)
		return "", err
	}
	return hash, nil
}

// GetBlock takes in a string of the hash and returns
// the block data as model.Block
func (c *Client) GetBlock(hash string) (*model.Block, error) {
	// Create a Resty Client
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	endpoint := fmt.Sprintf("%s/api/BlockStore/block", c.Endpoint)
	resp, err := c.api.R().
		SetQueryParams(map[string]string{
			"Hash":       hash,
			"OutputJson": "true",
		}).
		SetHeader("Accept", "application/json").
		Get(endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var block model.Block
	err = json.Unmarshal(resp.Body(), &block)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &block, nil
}

// GetTx takes in a string of the transaction id and
// returns the transaction data as model.Transaction
func (c *Client) GetTx(txid string) (*model.Transaction, error) {
	// Create a Resty Client
	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	endpoint := fmt.Sprintf("%s/api/Node/getrawtransaction", c.Endpoint)
	resp, err := c.api.R().
		SetQueryParams(map[string]string{
			"trxid":   txid,
			"verbose": "true",
		}).
		SetHeader("Accept", "application/json").
		Get(endpoint)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var transaction model.Transaction
	err = json.Unmarshal(resp.Body(), &transaction)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &transaction, nil
}

// GetVout returns the first vout that matches the scriptPubkey.type
// with the argument txType
// Returns an error if nothing is found
func GetVout(vout []model.Vout, txType string) (model.Vout, error) {
	for _, spendings := range vout {
		voutType := spendings.ScriptPubKey.Type
		if voutType == txType {
			return spendings, nil
		}
	}
	return model.Vout{}, fmt.Errorf("Vout not found")
}
