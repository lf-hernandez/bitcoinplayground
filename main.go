package main

import (
	"fmt"

	"github.com/lf-hernandez/bitcoinplayground/ledger"
	"github.com/lf-hernandez/bitcoinplayground/merkletree"
)

func main() {
	fmt.Print("Communal ledger with DSA signing simulation\n\n")
	ledger.RunSimulation()
	fmt.Println("\n=========================================================")
	fmt.Print("\nMerkle Tree simulation\n\n")
	merkletree.RunSimulation()
}
