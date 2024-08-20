package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/syndtr/goleveldb/leveldb"
)

type LogEntry struct {
	L1InfoRoot string
	BlockTime  time.Time
	ParentHash string
}


func initLevelDB() *leveldb.DB {
	db, err := leveldb.OpenFile("data/logs.db", nil)
	if err != nil {
		log.Fatalf("Failed to open LevelDB database: %v", err)
	}
	return db
}

func storeLogEntry(db *leveldb.DB, index int, entry LogEntry) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(entry); err != nil {
		log.Fatalf("Failed to encode log entry: %v", err)
	}
	err := db.Put([]byte(fmt.Sprintf("%d", index)), buf.Bytes(), nil)
	if err != nil {
		log.Fatalf("Failed to store log entry in LevelDB: %v", err)
	}
}

func main() {
	
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/replace with your project id to access the eth node in sepolia testnet")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	defer client.Close()

	
	db := initLevelDB()
	defer db.Close()

	contractAddress := common.HexToAddress("0x761d53b47334bee6612c0bd1467fb881435375b2")
	topic := common.HexToHash("0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766")

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{topic}},
	}


	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to get logs: %v", err)
	}

	
	index := 0
	for _, vLog := range logs {
		block, err := client.BlockByHash(context.Background(), vLog.BlockHash)
		if err != nil {
			log.Fatalf("Failed to get block: %v", err)
		}

		entry := LogEntry{
			L1InfoRoot: vLog.Topics[0].Hex(), 
			BlockTime:  time.Unix(int64(block.Time()), 0),
			ParentHash: block.ParentHash().Hex(),
		}

		storeLogEntry(db, index, entry)
		index++
	}
}
