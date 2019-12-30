package account

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/seaio-co/crypto/ecdsa"
	"github.com/seaio-co/util/serialize"
)

// Account can be managed
type Account struct {
	mu       sync.RWMutex
	path     string
	address  string
	keyStore string
}

// SetSavePath set a path to save keystore
func (a *Account) SetSavePath(path string) error {
	path, _ = filepath.Abs(path)
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return errors.New("not a dir")
	}
	a.mu.Lock()
	a.path = path
	a.mu.Unlock()
	return nil
}

// GetSavePath can return a path which the keystore saved
func (a *Account) GetSavePath() (path string, err error) {
	if a.path == "" {
		return "", errors.New("Path must be set first")
	}
	return a.path, nil
}

func (a *Account) saveKeyStore() error {
	filename := a.path + `/` + a.address
	a.mu.Lock()
	err := ioutil.WriteFile(filename, []byte(a.keyStore), 0644)
	a.mu.Unlock()
	return err
}

// CreateAccount 创建账户
func (a *Account) CreateAccount(password string) (address string, err error) {
	if a.path == "" {
		return "", errors.New("Path must be set first")
	}

	prv, pub := ecdsa.GenerateKey()

	keyJSON, err := generateKeyStore(prv, password)
	if err != nil {
		return "", err
	}

	a.address = strings.ToLower(pub.ToAddress().Hex())
	a.keyStore = keyJSON

	//保存keystore
	err = a.saveKeyStore()
	if err != nil {
		return "", err
	}
	return a.address, nil
}

// DelAccount will delete a Account data
func (a *Account) DelAccount(address, password string) (string, error) {
	if a.path == "" {
		return "", errors.New("Path must be set first")
	}
	filename := a.path + `/` + address
	// 读取keystore
	a.mu.RLock()
	keyJSON, err := ioutil.ReadFile(filename)
	a.mu.RUnlock()
	if err != nil {
		return "", err
	}
	keystore := string(keyJSON)
	// 验证密码
	_, err = decryptKeyStore(keystore, password)
	if err != nil {
		return "", err
	}

	a.mu.Lock()
	err = os.Remove(filename)
	a.mu.Unlock()
	if err != nil {
		return "", err
	}

	return address, nil
}

// EXportKeystore can export a keystore
func (a *Account) EXportKeystore(address, password string) (keystore string, err error) {
	if a.path == "" {
		return "", errors.New("Path must be set first")
	}
	filename := a.path + `/` + address

	a.mu.RLock()
	keyJSON, err := ioutil.ReadFile(filename)
	a.mu.RUnlock()
	if err != nil {
		return "", err
	}

	keystore = string(keyJSON)

	_, err = decryptKeyStore(keystore, password)
	if err != nil {
		return "", err
	}

	return
}

// ImportKeystore can import a keystroe
func (a *Account) ImportKeystore(keystore, password string) (address string, err error) {
	if a.path == "" {
		return "", errors.New("Path must be set first")
	}

	privatekey, err := decryptKeyStore(keystore, password)
	if err != nil {
		return "", err
	}
	address = strings.ToLower(privatekey.ToPubKey().ToAddress().Hex())

	// 如果用户手拼一个keystore文件，或者故意改写address为错误的地址，则使用正确地址保存keyStore
	keyJSON := new(KeyJSON)
	serialize.JsonUnMarshal([]byte(keystore), &keyJSON)
	if strings.Compare(keyJSON.Address, address) != 0 {
		keyJSON.Address = address
		k, _ := serialize.JsonMarshal(keyJSON)
		keystore = string(k)
	}

	a.address = address
	a.keyStore = keystore
	err = a.saveKeyStore()
	if err != nil {
		return "", err
	}

	return a.address, nil
}

// DecryptKeyStore 解密keystore
func (a *Account) DecryptKeyStore(keystore, password string) (*ecdsa.PrivateKey, error) {
	return decryptKeyStore(keystore, password)
}
