// Copyright (c) 2018 Clearmatics Technologies Ltd

package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

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
var blockHash = "0xdc69532823b281514c0eec2a809dfc6c690a6bbb3e9006df4914a5b69cbb6371"
var blockNumber = "2569550"

var endpoint = "http://127.0.0.1:9999"

func main() {
	txHash := common.HexToHash(txString)

	block := getBlockByTransaction(txHash)
	fmt.Printf("BlockHash:\n%+v\nReceived Blockhash:\n%+v\n\n", blockHash, block.BlockHash)

	// Connect to the EthClient
	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal("could not create RPC client: %v", err)
	}

	// _, _, err := client.TransactionByHash(context.Background(), txHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	blockNum := new(big.Int)
	blockNum.SetString(block.BlockNumber[2:], 16)

	// fmt.Println(block.BlockNumber[2:])
	// fmt.Printf(blockNum)
	fmt.Printf("BlockNumber:\n%v\nReceived BlockNumber:\n%v\n\n", blockNumber, blockNum)
	newBlock, err := client.BlockByNumber(context.Background(), blockNum)
	if err != nil {
		log.Fatal(err)
	}

	for idx, tx := range newBlock.Transactions() {
		fmt.Printf("Index:%v\t Transaction:%x\n", idx, tx)
	}
	// fmt.Printf("%+v\n", newBlock)
	// fmt.Printf("%+x\n", newBlock.Transactions)

}

func getBlockByTransaction(txHash common.Hash) (tx *transaction) {
	// Connect to the RPC Client
	client, err := rpc.Dial(endpoint)
	if err != nil {
		log.Fatalf("could not create RPC client: %v", err)
	}

	err = client.Call(&tx, "eth_getTransactionByHash", txHash)
	if err != nil {
		fmt.Println("can't get latest block:", err)
		return
	}

	return
}
