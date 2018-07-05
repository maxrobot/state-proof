package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/maxrobot/go-ethereum/ethdb"
	"github.com/maxrobot/go-ethereum/rlp"
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
		receipt, _ := client.TransactionReceipt(context.Background(), tx.Hash())
		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))   // rlp encode index of transaction
		rlpReceipt, _ := rlp.EncodeToBytes(receipt) // rlp encode receipt

		trieObj.Update(rlpIdx, rlpReceipt)

		// txRlpHash := crypto.Keccak256Hash(rlpReceipt)

		// fmt.Printf("TxHash[%d]: \t% 0x\n\tHash(RLP(Tx)): \t% 0x\n",
		// 	idx, tx.Hash().Bytes(), txRlpHash.Bytes())

		// Get the information about the receipt I care about...
		if transx == tx {
			txIdx = rlpIdx
			leaf = rlpReceipt
		}

	}

	root := trieObj.Hash()
	expectedRoot := block.ReceiptHash()

	fmt.Printf("\nExpected Root:\t%x\nRecovered Root:\t%x\n", expectedRoot, root)

	// Generate a merkle proof for a key
	proof := ethdb.NewMemDatabase()
	trieObj.Prove(txIdx, 0, proof)
	if proof == nil {
		fmt.Printf("prover: nil proof")
	}

	// Verify the proof
	val, _, err := trie.VerifyProof(root, txIdx, proof)
	if err != nil {
		fmt.Printf("prover: failed to verify proof: %v\nraw proof: %x", err, proof)
	}
	if !bytes.Equal(val, leaf) {
		fmt.Printf("prover: verified value mismatch: have %x, want 'k'", val)
	}
	fmt.Printf("\nVerified Value:\t%x\nExpected Leaf:\t%x\n", val, leaf)
}
