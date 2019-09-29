package net

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"
)

// Client
type Client struct {
	conn   net.Conn
	Server *server
}

// server
type server struct {
	address                  string
	config                   *tls.Config
	onNewClientCallback      func(c *Client)
	onClientConnectionClosed func(c *Client, err error)
	onNewMessage             func(c *Client, message string)
}

// listen
func (c *Client) listen() {
	c.Server.onNewClientCallback(c)
	reader := bufio.NewReader(c.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			c.conn.Close()
			c.Server.onClientConnectionClosed(c, err)
			return
		}
		c.Server.onNewMessage(c, message)
	}
}

// Send
func (c *Client) Send(message string) error {
	_, err := c.conn.Write([]byte(message))
	return err
}

// SendBytes
func (c *Client) SendBytes(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

// Conn
func (c *Client) Conn() net.Conn {
	return c.conn
}

// Close
func (c *Client) Close() error {
	return c.conn.Close()
}

// OnNewClient
func (s *server) OnNewClient(callback func(c *Client)) {
	s.onNewClientCallback = callback
}

// OnClientConnectionClosed
func (s *server) OnClientConnectionClosed(callback func(c *Client, err error)) {
	s.onClientConnectionClosed = callback
}

// OnNewMessage
func (s *server) OnNewMessage(callback func(c *Client, message string)) {
	s.onNewMessage = callback
}

// Listen
func (s *server) Listen() {
	var listener net.Listener
	var err error
	if s.config == nil {
		listener, err = net.Listen("tcp", s.address)
	} else {
		listener, err = tls.Listen("tcp", s.address, s.config)
	}
	if err != nil {
		log.Fatal("Error starting TCP server.")
	}
	defer listener.Close()

	for {
		conn, _ := listener.Accept()
		client := &Client{
			conn:   conn,
			Server: s,
		}
		go client.listen()
	}
}

// New
func New(address string) *server {
	log.Println("Creating server with address", address)
	server := &server{
		address: address,
		config:  nil,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// NewWithTLS
func NewWithTLS(address string, certFile string, keyFile string) *server {
	log.Println("Creating server with address", address)
	cert, _ := tls.LoadX509KeyPair(certFile, keyFile)
	config := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	server := &server{
		address: address,
		config:  &config,
	}

	server.OnNewClient(func(c *Client) {})
	server.OnNewMessage(func(c *Client, message string) {})
	server.OnClientConnectionClosed(func(c *Client, err error) {})

	return server
}

// CreateLink
func CreateLink(addr string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", addr, 3 * time.Second)
	if err != nil {
		return nil, err
	}
	return conn, nil
}