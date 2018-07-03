// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maxrobot/go-ethereum/ethclient"
)

type transaction struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Input            string `json:"input"`
}

var txString = "0xd0a72c0b6a703ad82853346c6b1ddaf75a9fc474100e4b9b301b42d8d88d9353"
var blockHash = "0xdc69532823b281514c0eec2a809dfc6c690a6bbb3e9006df4914a5b69cbb637"

func main() {
	// Connect to the RPC Client
	rpcClient, err := rpc.Dial("http://127.0.0.1:9999")
	if err != nil {
		log.Fatalf("could not create RPC client: %v", err)
	}

	txHash := common.HexToHash(txString)

	getTransactionProof(rpcClient, txHash)

	// Connect to the EthClient
	client, err := ethclient.Dial("http://127.0.0.1:9999")
	if err != nil {
		log.Fatal("could not create RPC client: %v", err)
	}

	tx, isPending, err := client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex())     // 0x5d49fcaa394c97ec8a9c3e7bd9e8388d420fb050a52083ca52ff24b3b65bc9c2
	fmt.Println(tx.Value().String()) // 10000000000000000 (in wei)
	fmt.Println(isPending)           // false
}

func getTransactionProof(client *rpc.Client, txHash common.Hash) {
	var tx *transaction
	err := client.Call(&tx, "eth_getTransactionByHash", txHash)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}
	fmt.Printf("Block Hash: %v\n", tx.BlockHash)

}
