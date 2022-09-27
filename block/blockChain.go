package block

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/boltdb/bolt"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
	lastBlockKey = "last"
)

type BlockChain struct {
	// Blocks []*Block
	tip []byte
	DB  *bolt.DB
}

var BC *BlockChain
var once sync.Once

func generateGenesis() *Block {
	return NewBlock("GENESIS BLOCK", []byte{})
}

func GetBlockChain() *BlockChain {
	var last []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket([]byte(blocksBucket))
		if bc == nil {
			genesis := generateGenesis()
			fmt.Println("Generate Genesis Block")
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				return err
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				return err
			}
			err = b.Put([]byte(lastBlockKey), genesis.Hash)
			if err != nil {
				return err
			}

			last = genesis.Hash
		} else {
			last = bc.Get([]byte(lastBlockKey))
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	bc := BlockChain{last, db}
	return &bc
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastBlockKey))
		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(lastBlockKey), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil
	})
}

func (bc *BlockChain) Iterator() *BlockChainInterator {
	bci := &BlockChainInterator{bc.tip, bc.DB}
	return bci
}

func (bc *BlockChain) ShowBlocks() {
	bci := bc.Iterator()

	for {
		block := bci.Next()
		pow := NewProofOfWork(block)

		fmt.Println("TimeStamp:", block.TimeStamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Prev Hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("is Validated: %s\n", strconv.FormatBool(pow.Validate()))

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
