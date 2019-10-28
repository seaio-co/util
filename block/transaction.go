package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10

// Transaction
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}


type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}


type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}