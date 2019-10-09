package net

import (
	"fmt"
)

type SetLingerer interface {
	SetLinger(sec int) error
}
type SetNoDelayer interface {
	SetNoDelay(noDelay bool) error
}

type SetReadBufferer interface {
	SetReadBuffer(bytes int) error
}
type SetWriteBufferer interface {
	SetWriteBuffer(bytes int) error
}
type SetReadBuffer interface {
	SetReadBuffer(bytes int) error
}

// SetNoDelay
func SetNoDelay(conn interface{}, noDelay bool) error {
	ccd, _ := conn.(SetNoDelayer)
	if ccd == nil {
		return fmt.Errorf("conn not provided SetNoDelay method")
	}
	ccd.SetNoDelay(noDelay)
	return nil
}

// SetLinger
func SetLinger(c interface{}, sec int) error {
	ccd, _ := c.(SetLingerer)
	if ccd == nil {
		return fmt.Errorf("conn not provided SetNoDelay method")
	}
	ccd.SetLinger(sec)
	return nil
}