package merkle

import (
	"bytes"
	"crypto/sha256"
	"errors"
)

var (
	DOUBLE_SHA256 = func(s []Uint256) Uint256 {
		b := new(bytes.Buffer)
		for _, d := range s {
			d.Serialize(b)
		}
		temp := sha256.Sum256(b.Bytes())
		f := sha256.Sum256(temp[:])
		return Uint256(f)
	}
)

type MerkleTree struct {
	Depth uint
	Root  *MerkleTreeNode
}

type MerkleTreeNode struct {
	Hash  Uint256
	Left  *MerkleTreeNode
	Right *MerkleTreeNode
}

func (t *MerkleTreeNode) IsLeaf() bool {
	return t.Left == nil && t.Right == nil
}

func NewMerkleTree(hashes []Uint256) (*MerkleTree, error) {
	if len(hashes) == 0 {
		return nil, errors.New("NewMerkleTree input no item error.")
	}
	var height uint

	height = 1
	nodes := generateLeaves(hashes)
	for len(nodes) > 1 {
		nodes = levelUp(nodes)
		height += 1
	}
	mt := &MerkleTree{
		Root:  nodes[0],
		Depth: height,
	}
	return mt, nil

}

func generateLeaves(hashes []Uint256) []*MerkleTreeNode {
	var leaves []*MerkleTreeNode
	for _, d := range hashes {
		node := &MerkleTreeNode{
			Hash: d,
		}
		leaves = append(leaves, node)
	}
	return leaves
}

func levelUp(nodes []*MerkleTreeNode) []*MerkleTreeNode {
	var nextLevel []*MerkleTreeNode
	for i := 0; i < len(nodes)/2; i++ {
		var data []Uint256
		data = append(data, nodes[i*2].Hash)
		data = append(data, nodes[i*2+1].Hash)
		hash := DOUBLE_SHA256(data)
		node := &MerkleTreeNode{
			Hash:  hash,
			Left:  nodes[i*2],
			Right: nodes[i*2+1],
		}
		nextLevel = append(nextLevel, node)
	}
	if len(nodes)%2 == 1 {
		var data []Uint256
		data = append(data, nodes[len(nodes)-1].Hash)
		data = append(data, nodes[len(nodes)-1].Hash)
		hash := DOUBLE_SHA256(data)
		node := &MerkleTreeNode{
			Hash:  hash,
			Left:  nodes[len(nodes)-1],
			Right: nodes[len(nodes)-1],
		}
		nextLevel = append(nextLevel, node)
	}
	return nextLevel
}

func ComputeRoot(hashes []Uint256) (Uint256, error) {
	if len(hashes) == 0 {
		return Uint256{}, errors.New("NewMerkleTree input no item error.")
	}
	if len(hashes) == 1 {
		return hashes[0], nil
	}
	tree, _ := NewMerkleTree(hashes)
	return tree.Root.Hash, nil
}

type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

func NewMerkleNode(left,right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	}else {
		prevHashes := append(left.Data,right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	mNode.Left = left
	mNode.Right = right

	return &mNode
}