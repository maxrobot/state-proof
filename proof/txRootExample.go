package main

import (
	"context"
	"encoding/hex"
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

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("we have a connection")

	// Select a specific block
	num := "5904064"
	blockNumber := new(big.Int)
	blockNumber.SetString(num, 10)

	header, err := client.HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	block, err := client.BlockByNumber(context.Background(), header.Number)
	if err != nil {
		log.Fatal(err)
	}

	// Select a transaction, index should be 49
	transx := block.Transaction(common.HexToHash("0xd828cadfcc7694d314058404506fc10a4dadac72aba68c286b34137da4156630"))
	var txIdx int
	fmt.Printf("\nTransaction:\n% 0x", transx.Hash().Bytes())

	// VERIFY TRANSACTION ROOT
	trieObj, _ := trie.New(common.Hash{}, trie.NewDatabase(ethdb.NewMemDatabase())) // empty trie
	for idx, tx := range block.Transactions() {

		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

		trieObj.Update(rlpIdx, rlpTransaction) // update trie with the rlp encode index and the rlp encoded transaction
		root, err := trieObj.Commit(nil)       // commit to database (which in this case is stored in memory)
		if err != nil {
			log.Fatalf("commit error: %v", err)
		}

		txRlpHash := crypto.Keccak256Hash(rlpTransaction)

		fmt.Printf("TxHash[%d]: \t% 0x\n\tHash(RLP(Tx)): \t% 0x\n\tTrieRoot: \t% 0x\n",
			idx, tx.Hash().Bytes(), txRlpHash.Bytes(), root.Bytes())

		// Find the index of the transaction above
		if transx == tx {
			txIdx = idx
		}
	}

	fmt.Printf("\n\nBlock number: %d \n\tBlock.TxHash:\t% 0x \n\tTrie.Root:\t% 0x\n",
		block.Number, block.TxHash().Bytes(), trieObj.Root())

	fmt.Printf("TxHash[%d]: \t% 0x\n", txIdx, transx.Hash().Bytes())

	// Get proof of the transaction
	key := transx.Hash().Bytes()
	proof := ethdb.NewMemDatabase()
	trieObj.Prove(key, 0, proof)

	rootString := hex.EncodeToString(trieObj.Root())
	root := common.HexToHash("0x" + rootString)
	fmt.Printf("\n%x\n", trieObj.Root())
	fmt.Printf("\n%x\n", root)

	// Validate that the proof generates the same trie root
	val, _, err := trie.VerifyProof(root, key, proof)
	fmt.Println(val)
	// fmt.Printf("Generated Trie Root:\n\t% 0x", val)
	// if err != nil {
	// 	t.Fatalf("prover %d: failed to verify proof for key %x: %v\nraw proof: %x", i, kv.k, err, proof)
	// }
	// if !bytes.Equal(val, kv.v) {
	// 	t.Fatalf("prover %d: verified value mismatch for key %x: have %x, want %x", i, kv.k, val, kv.v)
	// }

}
