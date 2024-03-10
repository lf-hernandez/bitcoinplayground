package merkletree

type MerkleNode struct {
	Hash []byte
	L    *MerkleNode
	R    *MerkleNode
}

type MerkleTree struct {
	Root *MerkleNode
}
