package main

import (
	// "github.com/lf-hernandez/bitcoinplayground/ledger"
	"crypto/sha256"
	"fmt"

	"github.com/lf-hernandez/bitcoinplayground/merkletree"
)

func main() {
	// ledger.RunSimulation()
	fmt.Print("Merkle Tree simulation\n\n")
	txs := []string{"Felipe sent Satoshi 1 coin", "Jane sent John 5 coins", "Satoshi sent Jane 3 coins", "John sent Felipe 7 coins", "Unbalanced*"}
	for _, tx := range txs {
		h := sha256.New()
		h.Write([]byte(tx))
		fmt.Printf("tx: %v\nhash: %x\n", tx, h.Sum(nil))
	}
	fmt.Println("")
	fmt.Println("Building Merkle Tree...")
	fmt.Println("Unbalanced*")
	ubroot := merkletree.BuildTree(txs, true)
	ubtree := &merkletree.MerkleTree{Root: ubroot}
	ubtree.PrintTree()
	fmt.Println("Balanced")
	root := merkletree.BuildTree(txs, false)
	tree := &merkletree.MerkleTree{Root: root}
	tree.PrintTree()
}
