package account

import (
	"bytes"
	"encoding/hex"
	"errors"
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

func generateKeyStore(prvKey *ecdsa.PrivateKey, password string) (string, error) {

	salt := common.GetEntropyCSPRNG(32)
	derivedKey, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return "", err
	}

	encryptKey := derivedKey[:32]
	keyBytes := math.PaddedBigBytes(prvKey.D, 32)

	iv := common.GetEntropyCSPRNG(16)
	cipherText, err := aes.AesCTRXOR(encryptKey, keyBytes, iv)
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

func decryptKeyStore(keystore, password string) (*ecdsa.PrivateKey, error) {

	keyJson := new(KeyJSON)
	err := serialize.JsonUnMarshal([]byte(keystore), &keyJson)
	if err != nil {
		return nil, err
	}

	salt, err := hex.DecodeString(keyJson.ScryptParams.Salt)
	if err != nil {
		return nil, err
	}
	derivedKey, err := scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptDKLen)
	if err != nil {
		return nil, err
	}
	cipherText, err := hex.DecodeString(keyJson.ScryptParams.Text)
	if err != nil {
		return nil, err
	}
	calculatedMAC := sha3.Keccak256(derivedKey[16:32], cipherText) //验证时输入密码获得的mac
	mac, err := hex.DecodeString(keyJson.ScryptParams.Mac)         //生成keystore时输入密码获得的mac
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(calculatedMAC, mac) { //两次mac相等，则证明密码正确
		return nil, errors.New("could not decrypt key with given passphrase")
	}
	// 使用验证过的、正确的密码把cipherText还原为“原文”，并得到私钥
	iv, err := hex.DecodeString(keyJson.ScryptParams.IV)
	if err != nil {
		return nil, err
	}
	encryptKey := derivedKey[:32]
	plainText, err := aes.AesCTRXOR(encryptKey, cipherText, iv)
	if err != nil {
		return nil, err
	}

	prvKey, _ := ecdsa.ToECDSA(plainText, true)
	if err != nil {
		return nil, err
	}
	return prvKey, nil
}
