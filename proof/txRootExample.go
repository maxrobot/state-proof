package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/maxrobot/go-ethereum/crypto"
)

var expectedBlockHash = common.HexToHash("0xa68e48252380b056f5e1c1a738897b84148be4ffbcf0857effbaa086a8a99fbb")
var expectedTx = common.HexToHash("0xd828cadfcc7694d314058404506fc10a4dadac72aba68c286b34137da4156630")

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to: https://mainnet.infura.io")

	// Select a specific block
	num := "5904064"
	blockNumber := new(big.Int)
	blockNumber.SetString(num, 10)

	// Fetch header of block num
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	// assert.Equal(t, expectedBlockHash, header.Hash())

	// Fetch block of block num
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}
	// assert.Equal(t, expectedBlockHash, block.Hash())

	// Select a transaction, index should be 49
	transx := block.Transaction(expectedTx)
	var txIdx []byte
	var leaf []byte
	// fmt.Printf("\nTransaction:\n% 0x", transx.Hash().Bytes())

	// Generate the trie
	trieObj := new(trie.Trie)
	for idx, tx := range block.Transactions() {
		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

		trieObj.Update(rlpIdx, rlpTransaction)

		txRlpHash := crypto.Keccak256Hash(rlpTransaction)

		fmt.Printf("TxHash[%d]: \t% 0x\n\tHash(RLP(Tx)): \t% 0x\n",
			idx, tx.Hash().Bytes(), txRlpHash.Bytes())

		// Get the information about the transaction I care about...
		if transx == tx {
			txIdx = rlpIdx
			leaf = rlpTransaction
		}

	}

	// Generate a proof
	proof := ethdb.NewMemDatabase()
	trieObj.Prove(txIdx, 0, proof)
	if proof == nil {
		fmt.Printf("prover: nil proof")
	}

	// Verify the proof
	val, _, err := trie.VerifyProof(trieObj.Hash(), txIdx, proof)
	if err != nil {
		fmt.Printf("prover: failed to verify proof: %v\nraw proof: %x", err, proof)
	}
	if !bytes.Equal(val, leaf) {
		fmt.Printf("prover: verified value mismatch: have %x, want 'k'", val)
	}
	fmt.Printf("\nVerified Value:\t%x\nExpected Leaf:\t%x\n", val, leaf)
}
