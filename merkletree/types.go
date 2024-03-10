package merkletree

type MerkleNode struct {
	hash []byte
	l    *MerkleNode
	r    *MerkleNode
}

type MerkleTree struct {
	Root *MerkleNode
}
