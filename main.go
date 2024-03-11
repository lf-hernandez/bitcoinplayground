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

	txs := []string{
		"Felipe sent Satoshi 1 coin",
		"Jane sent John 5 coins",
		"Satoshi sent Jane 3 coins",
		"John sent Felipe 7 coins",
	}
	for _, tx := range txs {
		h := sha256.New()
		h.Write([]byte(tx))
		fmt.Printf("tx: %v\nhash: %x\n", tx, h.Sum(nil))
	}
	fmt.Println("\nBuilding Merkle tree")
	root := merkletree.BuildTree(txs, false)
	tree := &merkletree.MerkleTree{Root: root}
	tree.PrintTree()

	fmt.Println("\nGenerating Merkle proof")
	proof := tree.GenerateMerkleProof(txs[3])
	for _, e := range proof {
		fmt.Printf("%x\n", e.Hash)
	}

	h := sha256.New()
	h.Write([]byte(txs[3]))
	targetHash := h.Sum(nil)

	fmt.Printf(
		"\nRunning verification\nis transaction: %v\nwith hash: %x \nin current tree: %v\n",
		txs[3],
		targetHash,
		merkletree.VerifyMerkleProof(targetHash, proof, tree.Root.Hash),
	)

}
