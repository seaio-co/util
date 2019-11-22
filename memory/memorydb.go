package memory

import (
	"errors"
	"sync"
)

type MemDB struct {
	db   map[string][]byte
	lock sync.RWMutex
}

var mdb *MemDB
var err error
var once sync.Once

// Init
func Init() *MemDB {
	once.Do(func() {
		mdb, err = newMemDB()
		if err != nil {
			panic(err)
		}
	})
	return mdb
}

// newMemDB
func newMemDB() (*MemDB, error) {
	return &MemDB{
		db: make(map[string][]byte),
	}, nil
}

// Put
func (db *MemDB) Put(key []byte, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.db[string(key)] = value
	return nil
}

// Get
func (db *MemDB) Get(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if value, ok := db.db[string(key)]; ok {
		return value, nil
	}

	return nil, errors.New("key is not found")
}

// Has
func (db *MemDB) Has(key []byte) (bool, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if _, ok := db.db[string(key)]; ok {
		return true, nil
	}
	return false, nil
}

// Del
func (db *MemDB) Del(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	delete(db.db, string(key))

	return nil
}

// GetAllKey
func (db *MemDB) GetAllKey() [][]byte {
	db.lock.RLock()
	defer db.lock.RUnlock()

	keys := [][]byte{}
	for key := range db.db {
		keys = append(keys, []byte(key))
	}
	return keys
}

// Path
func (db *MemDB) Path() string {
	return ""
}

// Close
func (db *MemDB) Close() error {
	return nil
}