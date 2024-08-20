package data

import (
	"bytes"
	"encoding/gob"

	"github.com/syndtr/goleveldb/leveldb"
)

type LogEntry struct {
	L1InfoRoot string
	BlockTime  int64
	ParentHash string
}


func Serialize(entry LogEntry) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(entry)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}


func Deserialize(data []byte) (LogEntry, error) {
	var entry LogEntry
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&entry)
	if err != nil {
		return LogEntry{}, err
	}
	return entry, nil
}


func StoreLogEntry(db *leveldb.DB, index int, entry LogEntry) error {
	data, err := Serialize(entry)
	if err != nil {
		return err
	}
	err = db.Put([]byte(indexKey(index)), data, nil)
	if err != nil {
		return err
	}
	return nil
}


func RetrieveLogEntry(db *leveldb.DB, index int) (LogEntry, error) {
	data, err := db.Get([]byte(indexKey(index)), nil)
	if err != nil {
		return LogEntry{}, err
	}
	return Deserialize(data)
}

func indexKey(index int) string {
	return string(index)
}
