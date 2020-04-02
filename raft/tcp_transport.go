package raft

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	errNotAdvertisable = errors.New("local bind address is not advertisable")
	errNotTCP          = errors.New("local address is not a TCP address")
)

// TCPStreamLayer implements StreamLayer interface for plain TCP.
type TCPStreamLayer struct {
	advertise net.Addr
	listener  *net.TCPListener
}

// NewTCPTransport returns a NetworkTransport that is built on top of
// a TCP streaming transport layer.
func NewTCPTransport(
	bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	logOutput io.Writer,
) (*NetworkTransport, error) {
	return newTCPTransport(bindAddr, advertise, func(stream StreamLayer) *NetworkTransport {
		return NewNetworkTransport(stream, maxPool, timeout, logOutput)
	})
}

// NewTCPTransportWithLogger returns a NetworkTransport that is built on top of
// a TCP streaming transport layer, with log output going to the supplied Logger
func NewTCPTransportWithLogger(
	bindAddr string,
	advertise net.Addr,
	maxPool int,
	timeout time.Duration,
	logger hclog.Logger,
) (*NetworkTransport, error) {
	return newTCPTransport(bindAddr, advertise, func(stream StreamLayer) *NetworkTransport {
		return NewNetworkTransportWithLogger(stream, maxPool, timeout, logger)
	})
}

// NewTCPTransportWithConfig returns a NetworkTransport that is built on top of
// a TCP streaming transport layer, using the given config struct.
func NewTCPTransportWithConfig(
	bindAddr string,
	advertise net.Addr,
	config *NetworkTransportConfig,
) (*NetworkTransport, error) {
	return newTCPTransport(bindAddr, advertise, func(stream StreamLayer) *NetworkTransport {
		config.Stream = stream
		return NewNetworkTransportWithConfig(config)
	})
}
