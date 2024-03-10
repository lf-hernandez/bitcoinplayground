package merkletree

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func BuildTree(txs []string) *MerkleNode {
	var leaves []*MerkleNode

	for _, tx := range txs {
		h := sha256.New()
		h.Write([]byte(tx))
		hashedContent := h.Sum(nil)
		leaves = append(leaves, &MerkleNode{hash: hashedContent})
	}

	return buildTree(leaves)
}

func buildTree(nodes []*MerkleNode) *MerkleNode {
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

	return buildTree(higherLevelNodes)
}

func composeNodes(l, r *MerkleNode) *MerkleNode {
	h := sha256.New()
	h.Write(l.hash)
	h.Write(r.hash)
	compositeHash := h.Sum(nil)
	return &MerkleNode{hash: compositeHash, l: l, r: r}
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

	fmt.Printf("%s\n", hex.EncodeToString(n.hash))

	if n.l != nil || n.r != nil {
		fmt.Printf("%sLeft: ", indent)
		printNode(n.l, level+1)
		fmt.Printf("%sRight: ", indent)
		printNode(n.r, level+1)
	}
}
