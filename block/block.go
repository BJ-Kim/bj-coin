package block

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	TimeStamp     int64  `validate:"required"`
	Data          []byte `validate:"required"`
	PrevBlockHash []byte `validate:"required"`
	Hash          []byte `validate:"required"`
	Nonce         int    `validate:"min=0"`
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	newBlock := &Block{time.Now().Unix(), []byte(data), prevBlockHash, nil, 0}
	// block.SetHash()
	pow := NewProofOfWork(newBlock)
	nonce, hash := pow.Run()
	newBlock.Hash = hash[:]
	newBlock.Nonce = nonce

	return newBlock
}

func NewGenesisBlock() *Block {
	return NewBlock("GENESIS BLOCK", []byte{})
}

func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}
	return &block
}
