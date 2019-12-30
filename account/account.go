package account

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/errors"
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
