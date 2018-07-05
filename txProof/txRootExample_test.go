package main

import (
	"bytes"
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/maxrobot/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
)

var expectedBlockHash = common.HexToHash("0xa68e48252380b056f5e1c1a738897b84148be4ffbcf0857effbaa086a8a99fbb")
var expectedTx = common.HexToHash("0xd828cadfcc7694d314058404506fc10a4dadac72aba68c286b34137da4156630")

func updateString(trie *trie.Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}

func TestOneElementProof(t *testing.T) {
	trieObj := new(trie.Trie)
	var idx = "spoon"
	var leaf = "kita"

	// Put some stuff in the trie
	updateString(trieObj, "niki", "fmg")
	updateString(trieObj, "spoon", "kita")
	updateString(trieObj, "spoon2", "kita2")
	updateString(trieObj, "spoon3", "kita3")
	updateString(trieObj, "spoon4", "kita4")

	// Generate a proof
	proof := ethdb.NewMemDatabase()
	trieObj.Prove([]byte(idx), 0, proof)
	if proof == nil {
		t.Fatalf("prover: nil proof")
	}

	// Verify the proof
	val, _, err := trie.VerifyProof(trieObj.Hash(), []byte(idx), proof)
	if err != nil {
		t.Fatalf("prover: failed to verify proof: %v\nraw proof: %x", err, proof)
	}
	if !bytes.Equal(val, []byte(leaf)) {
		t.Fatalf("prover: verified value mismatch: have %x, want 'k'", val)
	}

}

func TestOneTransactionProof(t *testing.T) {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	// Select a specific block
	num := "5904064"
	blockNumber := new(big.Int)
	blockNumber.SetString(num, 10)

	// Fetch header of block num
	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedBlockHash, header.Hash())

	// Fetch block of block num
	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expectedBlockHash, block.Hash())

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
		// fmt.Printf("%x\t%x\n", transx.Hash().Bytes(), tx.Hash().Bytes())
		if transx == tx {
			txIdx = rlpIdx
			leaf = rlpTransaction
		}

	}

	// Generate a proof
	proof := ethdb.NewMemDatabase()
	trieObj.Prove(txIdx, 0, proof)
	if proof == nil {
		t.Fatalf("prover: nil proof")
	}

	// Verify the proof
	val, _, err := trie.VerifyProof(trieObj.Hash(), txIdx, proof)
	if err != nil {
		t.Fatalf("prover: failed to verify proof: %v\nraw proof: %x", err, proof)
	}
	if !bytes.Equal(val, leaf) {
		t.Fatalf("prover: verified value mismatch: have %x, want 'k'", val)
	}
	// fmt.Printf("\nVerified Value:\t%x\nExpected Leaf:\t%x\n", val, leaf)

}
