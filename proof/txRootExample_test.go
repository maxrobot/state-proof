package main

import (
	"bytes"
	crand "crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/maxrobot/go-ethereum/crypto"
	"github.com/maxrobot/go-ethereum/trie"
)

type kv struct {
	k, v []byte
	t    bool
}

var expectedBlockHash = common.HexToHash("0xa68e48252380b056f5e1c1a738897b84148be4ffbcf0857effbaa086a8a99fbb")
var expectedTx = common.HexToHash("0xd828cadfcc7694d314058404506fc10a4dadac72aba68c286b34137da4156630")

// makeProvers creates Merkle trie provers based on different implementations to
// test all variations.
func makeProvers(trieObj *trie.Trie) []func(key []byte) *ethdb.MemDatabase {
	var provers []func(key []byte) *ethdb.MemDatabase

	// Create a direct trie based Merkle prover
	provers = append(provers, func(key []byte) *ethdb.MemDatabase {
		proof := ethdb.NewMemDatabase()
		trieObj.Prove(key, 0, proof)
		return proof
	})
	// Create a leaf iterator based Merkle prover
	provers = append(provers, func(key []byte) *ethdb.MemDatabase {
		proof := ethdb.NewMemDatabase()
		if it := trie.NewIterator(trieObj.NodeIterator(key)); it.Next() && bytes.Equal(key, it.Key) {
			for _, p := range it.Prove() {
				proof.Put(crypto.Keccak256(p), p)
			}
		}
		return proof
	})
	return provers
}

func randomTrie(n int) (*trie.Trie, map[string]*kv) {
	trieObj := new(trie.Trie)
	vals := make(map[string]*kv)
	for i := byte(0); i < 100; i++ {
		value := &kv{common.LeftPadBytes([]byte{i}, 32), []byte{i}, false}
		value2 := &kv{common.LeftPadBytes([]byte{i + 10}, 32), []byte{i}, false}
		trieObj.Update(value.k, value.v)
		trieObj.Update(value2.k, value2.v)
		vals[string(value.k)] = value
		vals[string(value2.k)] = value2
	}
	for i := 0; i < n; i++ {
		value := &kv{randBytes(32), randBytes(20), false}
		trieObj.Update(value.k, value.v)
		vals[string(value.k)] = value
	}
	return trieObj, vals
}

func randBytes(n int) []byte {
	r := make([]byte, n)
	crand.Read(r)
	return r
}

func TestProofMod(t *testing.T) {
	trieObj, vals := randomTrie(500)
	root := trieObj.Hash()
	for i, prover := range makeProvers(trieObj) {
		for _, kv := range vals {
			proof := prover(kv.k)
			if proof == nil {
				t.Fatalf("prover %d: missing key %x while constructing proof", i, kv.k)
			}
			val, _, err := trie.VerifyProof(root, kv.k, proof)
			if err != nil {
				t.Fatalf("prover %d: failed to verify proof for key %x: %v\nraw proof: %x", i, kv.k, err, proof)
			}
			if !bytes.Equal(val, kv.v) {
				t.Fatalf("prover %d: verified value mismatch for key %x: have %x, want %x", i, kv.k, val, kv.v)
			}
		}
	}
}

// func TestProof(t *testing.T) {
// 	trieObj, vals := randomTrie(500)
// 	root := trieObj.Hash()
// 	for i, prover := range makeProvers(trieObj) {
// 		for _, kv := range vals {
// 			proof := prover(kv.k)
// 			if proof == nil {
// 				t.Fatalf("prover %d: missing key %x while constructing proof", i, kv.k)
// 			}
// 			val, _, err := trie.VerifyProof(root, kv.k, proof)
// 			if err != nil {
// 				t.Fatalf("prover %d: failed to verify proof for key %x: %v\nraw proof: %x", i, kv.k, err, proof)
// 			}
// 			if !bytes.Equal(val, kv.v) {
// 				t.Fatalf("prover %d: verified value mismatch for key %x: have %x, want %x", i, kv.k, val, kv.v)
// 			}
// 		}
// 	}
// }

// func Test_txRootExample(t *testing.T) {
// 	client, err := ethclient.Dial("https://mainnet.infura.io")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("we have a connection")

// 	// Select a specific block
// 	num := "5904064"
// 	blockNumber := new(big.Int)
// 	blockNumber.SetString(num, 10)

// 	// Fetch header of block num
// 	header, err := client.HeaderByNumber(context.Background(), blockNumber)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	assert.Equal(t, expectedBlockHash, header.Hash())

// 	// Fetch block of block num
// 	block, err := client.BlockByNumber(context.Background(), header.Number)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	assert.Equal(t, expectedBlockHash, block.Hash())

// 	// Select a transaction, index should be 49
// 	transx := block.Transaction(expectedTx)
// 	var txIdx []byte
// 	fmt.Printf("\nTransaction:\n% 0x", transx.Hash().Bytes())

// 	// Generate the trie
// 	trieObj := new(trie.Trie)
// 	for idx, tx := range block.Transactions() {
// 		rlpIdx, _ := rlp.EncodeToBytes(uint(idx))  // rlp encode index of transaction
// 		rlpTransaction, _ := rlp.EncodeToBytes(tx) // rlp encode transaction

// 		trieObj.Update(rlpIdx, rlpTransaction)
// 		fmt.Printf("%x\t%x\n", transx.Hash().Bytes(), tx.Hash().Bytes())
// 		if transx == tx {
// 			fmt.Println(rlpIdx)
// 			txIdx = rlpIdx
// 		}

// 	}

// 	// Root hash
// 	root := trieObj.Hash()
// 	fmt.Printf("\nTransaction:\n% 0x\n", root)

// 	// Create the proof for transaction
// 	key, _ := rlp.EncodeToBytes(transx)
// 	proof := ethdb.NewMemDatabase()
// 	trieObj.Prove(key, 0, proof)
// 	if proof == nil {
// 		fmt.Println("Error no proof produced!")
// 	}

// 	fmt.Println(txIdx)
// 	// Verify the proof exists
// 	val, _, err := trie.VerifyProof(root, txIdx, proof)
// 	fmt.Println(val)

// }
