package account

import (
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
	"github.com/seaio-co/crypto/aes"
	"github.com/seaio-co/crypto/ecdsa"
	"github.com/seaio-co/crypto/sha3"
	"github.com/seaio-co/util/common"
	"github.com/seaio-co/util/math"
	"github.com/seaio-co/util/serialize"
	"golang.org/x/crypto/scrypt"
)

// KeyJSON keyStore结构体
type KeyJSON struct {
	ID           string        `json:"id"`
	Address      string        `json:"address"`
	Version      string        `json:"version"`
	ScryptParams *ScryptParams `json:"scryptParams"`
}

type ScryptParams struct {
	Salt string `json:"salt"`
	Mac  string `json:"mac"`
	IV   string `json:"ctr"`
	Text string `json:"text"`
}

const (
	version     = "1.0"
	scryptN     = 2048
	scryptP     = 6
	scryptR     = 8
	scryptDKLen = 32
)

// encryptKey 生成keystore算法
func generateKeyStore(prvKey *ecdsa.PrivateKey, password string) (string, error) {

	salt := common.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return "", err
	}

	encryptKey := derivedKey[:32]
	keyBytes := math.PaddedBigBytes(prvKey.D, 32)

	iv := common.GetEntropyCSPRNG(scryptDKLen)
	cipherText, err := aes.AesCTRXOR(encryptKey, keyBytes, iv[:aes.BlockSise])
	if err != nil {
		return "", err
	}

	mac := sha3.Keccak256(derivedKey[16:32], cipherText)
	scryptParams := new(ScryptParams)
	scryptParams.IV = hex.EncodeToString(iv)
	scryptParams.Mac = hex.EncodeToString(mac)
	scryptParams.Salt = hex.EncodeToString(salt)
	scryptParams.Text = hex.EncodeToString(cipherText)
	randomId, _ := uuid.NewUUID()

	keyJson := &KeyJSON{
		ID:           randomId.String(),
		Address:      strings.ToLower(prvKey.ToPubKey().ToAddress().Hex()),
		Version:      version,
		ScryptParams: scryptParams,
	}
	jsonByte, err := serialize.JsonMarshal(keyJson)
	if err != nil {
		return "", err
	}
	return string(jsonByte), nil
}
