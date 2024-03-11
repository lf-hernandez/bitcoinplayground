package main

import (
	// "github.com/lf-hernandez/bitcoinplayground/ledger"
	"crypto/sha256"
	"encoding/hex"
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

	// FIXME: unrefined approach
	fmt.Println("Unfrefined proof generation")
	var merkleProof []*merkletree.MerkleNode

	fmt.Printf("Verifying tx3: %v\n", txs[2])
	h := sha256.New()
	h.Write([]byte(txs[2]))
	hash := h.Sum(nil)

	tx3Leaf := merkletree.FindLeafByHash(hash, tree.Root)
	fmt.Printf("Found leaf: %x\n", tx3Leaf.Hash)
	merkleProof = append(merkleProof, tx3Leaf)

	h = sha256.New()
	h.Write([]byte(txs[3]))
	merkleProof = append(merkleProof, &merkletree.MerkleNode{Hash: h.Sum(nil)})

	h = sha256.New()
	h.Write([]byte(txs[0]))
	l1 := &merkletree.MerkleNode{Hash: h.Sum(nil)}
	h = sha256.New()
	h.Write([]byte(txs[1]))
	l2 := &merkletree.MerkleNode{Hash: h.Sum(nil)}

	h = sha256.New()
	h.Write(l1.Hash)
	h.Write(l2.Hash)
	chl1l2 := &merkletree.MerkleNode{Hash: h.Sum(nil)}

	merkleProof = append(merkleProof, chl1l2)
	merkleProof = append(merkleProof, root)

	fmt.Println("Building Merkle proof")
	for _, mn := range merkleProof {
		fmt.Println(hex.EncodeToString(mn.Hash))
	}

}
