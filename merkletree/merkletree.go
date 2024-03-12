package merkletree

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func VerifyMerkleProof(txHash []byte, proof []MerkleProofElement, rootHash []byte) bool {
	currentHash := txHash

	for _, e := range proof {
		if e.isLeft {
			currentHash = composeHash(e.Hash, currentHash)
		} else {
			currentHash = composeHash(currentHash, e.Hash)
		}
	}
	return bytes.Equal(currentHash, rootHash)
}

func composeHash(leftHash, rootHash []byte) []byte {
	h := sha256.New()
	h.Write(leftHash)
	h.Write(rootHash)
	return h.Sum(nil)
}

func (mt *MerkleTree) GenerateMerkleProof(tx string) []MerkleProofElement {
	h := sha256.New()
	h.Write([]byte(tx))
	txH := h.Sum(nil)
	leaf := mt.FindLeafByHash(txH)

	proof := mt.buildProof(leaf)
	return proof
}

func (mt *MerkleTree) buildProof(node *MerkleNode) []MerkleProofElement {
	var proof []MerkleProofElement

	for node != nil && node != mt.Root {
		sibling := findSibling(node)
		if sibling != nil {
			e := MerkleProofElement{
				Hash:    sibling.Hash,
				isLeft:  sibling == node.Parent.L,
				isRight: sibling == node.Parent.R,
			}
			proof = append(proof, e)
		}
		node = node.Parent
	}
	return proof
}

func findSibling(node *MerkleNode) *MerkleNode {
	p := node.Parent
	if p == nil {
		return nil
	}

	if node == p.L {
		return p.R
	}
	return p.L
}

func (mt *MerkleTree) FindLeafByHash(hash []byte) *MerkleNode {
	return findLeafByHash(hash, mt.Root)
}

func findLeafByHash(hash []byte, node *MerkleNode) *MerkleNode {
	if node == nil {
		return nil
	}

	if node.L == nil && node.R == nil {
		if bytes.Equal(node.Hash, hash) {
			return node
		}
	}

	leftLeaf := findLeafByHash(hash, node.L)

	if leftLeaf != nil {
		return leftLeaf
	}

	rightLeaf := findLeafByHash(hash, node.R)
	return rightLeaf
}

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

	return buildBalancedTree(leaves)
}

func buildUnbalancedTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var higherLevelNodes []*MerkleNode

	for i := 0; i < len(nodes); i += 2 {
		if i+1 < len(nodes) {
			compositeNode := composeNodes(nodes[i], nodes[i+1])
			higherLevelNodes = append(higherLevelNodes, compositeNode)
		} else {
			higherLevelNodes = append(higherLevelNodes, nodes[i])
		}
	}

	return buildUnbalancedTree(higherLevelNodes)
}

func buildBalancedTree(nodes []*MerkleNode) *MerkleNode {
	if len(nodes) == 1 {
		return nodes[0]
	}

	var hlns []*MerkleNode

	if len(nodes)%2 != 0 {
		ln := nodes[len(nodes)-1]
		nodes = append(nodes, ln)
	}

	for i := 0; i < len(nodes); i += 2 {
		cn := composeNodes(nodes[i], nodes[i+1])
		hlns = append(hlns, cn)
	}

	return buildBalancedTree(hlns)

}

func composeNodes(leftNode, rightNode *MerkleNode) *MerkleNode {
	h := sha256.New()
	h.Write(leftNode.Hash)
	h.Write(rightNode.Hash)
	compositeHash := h.Sum(nil)
	compositeNode := &MerkleNode{Hash: compositeHash, L: leftNode, R: rightNode}
	leftNode.Parent, rightNode.Parent = compositeNode, compositeNode
	return compositeNode
}

func (mt *MerkleTree) PrintTree() {
	fmt.Print("Root: ")
	printNode(mt.Root, 0)
}

func printNode(node *MerkleNode, level int) {
	if node == nil {
		return
	}

	indent := ""
	for i := 0; i < level; i++ {
		indent += "\t"
	}

	fmt.Printf("%s\n", hex.EncodeToString(node.Hash))

	if node.L != nil || node.R != nil {
		fmt.Printf("%sLeft: ", indent)
		printNode(node.L, level+1)
		fmt.Printf("%sRight: ", indent)
		printNode(node.R, level+1)
	}
}

func RunSimulation() {
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
	root := BuildTree(txs, false)
	tree := &MerkleTree{Root: root}
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
		VerifyMerkleProof(targetHash, proof, tree.Root.Hash),
	)
}
