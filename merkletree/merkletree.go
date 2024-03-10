package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func BuildTree(txs []string, isUnbalanced bool) *MerkleNode {
	var leaves []*MerkleNode

	for _, tx := range txs {
		h := sha256.New()
		h.Write([]byte(tx))
		hashedContent := h.Sum(nil)
		leaves = append(leaves, &MerkleNode{Hash: hashedContent})
	}

	if isUnbalanced {
		return buildUnbalancedTree(leaves)
	}

	return buildTree(leaves)
}

func buildUnbalancedTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var higherLevelNodes []*MerkleNode

	for i := 0; i < len(nodes); i += 2 {
		if i+1 < len(nodes) {
			compositeNodes := composeNodes(nodes[i], nodes[i+1])
			higherLevelNodes = append(higherLevelNodes, compositeNodes)
		} else {
			higherLevelNodes = append(higherLevelNodes, nodes[i])
		}
	}

	return buildUnbalancedTree(higherLevelNodes)
}

func buildTree(ns []*MerkleNode) *MerkleNode {
	if len(ns) == 1 {
		return ns[0]
	}

	var hlns []*MerkleNode

	if len(ns)%2 != 0 {
		ln := ns[len(ns)-1]
		ns = append(ns, ln)
	}

	for i := 0; i < len(ns); i += 2 {
		cn := composeNodes(ns[i], ns[i+1])
		hlns = append(hlns, cn)
	}

	return buildTree(hlns)

}

func composeNodes(l, r *MerkleNode) *MerkleNode {
	h := sha256.New()
	h.Write(l.Hash)
	h.Write(r.Hash)
	compositeHash := h.Sum(nil)
	return &MerkleNode{Hash: compositeHash, L: l, R: r}
}

func (m *MerkleTree) PrintTree() {
	fmt.Print("Merkle Root: ")
	printNode(m.Root, 0)
}

func printNode(n *MerkleNode, level int) {
	if n == nil {
		return
	}

	indent := ""
	for i := 0; i < level; i++ {
		indent += "\t"
	}

	fmt.Printf("%s\n", hex.EncodeToString(n.Hash))

	if n.L != nil || n.R != nil {
		fmt.Printf("%sLeft: ", indent)
		printNode(n.L, level+1)
		fmt.Printf("%sRight: ", indent)
		printNode(n.R, level+1)
	}
}
