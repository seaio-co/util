package block

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	Version      int64
	PreBlockHash []byte
	Hash         []byte
	TimeStamp    int64
	TargetBits   int64
	Nonce int64
	MerkelRoot []byte
	Data []byte
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	//initial block data
	block := &Block{
		Version:1,
		PreBlockHash:prevBlockHash,
		//Hash:[]
		TimeStamp:time.Now().Unix(),
		TargetBits:10,
		Nonce:5,
		MerkelRoot:[]byte{},
		Data:[]byte(data),
	}
	block.SetHash() //get block hash
	return block
}

func (block *Block)SetHash()  {
	tmp := [][]byte{
		IntToByte(block.Version),
		block.PreBlockHash,
		IntToByte(block.TimeStamp),
		block.MerkelRoot,
		IntToByte(block.Nonce),
		block.Data,
	}
	data := bytes.Join(tmp, []byte{})
	hash := sha256.Sum256(data)
	block.Hash = hash[:]
}

package main

import "os"

type BlockChain struct {
	blocks []*Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGensisBlock()}}
}

func (bc *BlockChain)AddBlock(data string)  {
	if len(bc.blocks) <= 0{
		os.Exit(1)
	}
	lastBlock := bc.blocks[len(bc.blocks)-1]
	prevBlockHash := lastBlock.Hash
	block := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, block)
}