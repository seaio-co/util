package EventBus

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type NetworkBus struct {
	*Client
	*Server
	service   *NetworkBusService
	sharedBus Bus
	address   string
	path      string
}

func NewNetworkBus(address, path string) *NetworkBus {
	bus := new(NetworkBus)
	bus.sharedBus = New()
	bus.Server = NewServer(address, path, bus.sharedBus)
	bus.Client = NewClient(address, path, bus.sharedBus)
	bus.service = &NetworkBusService{&sync.WaitGroup{}, false}
	bus.address = address
	bus.path = path
	return bus
}

func (networkBus *NetworkBus) EventBus() Bus {
	return networkBus.sharedBus
}

type NetworkBusService struct {
	wg      *sync.WaitGroup
	started bool
}

func (networkBus *NetworkBus) Start() error {
	var err error
	service := networkBus.service
	clientService := networkBus.Client.service
	serverService := networkBus.Server.service
	if !service.started {
		server := rpc.NewServer()
		server.RegisterName("ServerService", serverService)
		server.RegisterName("ClientService", clientService)
		server.HandleHTTP(networkBus.path, "/debug"+networkBus.path)
		l, e := net.Listen("tcp", networkBus.address)
		if e != nil {
			err = fmt.Errorf("listen error: %v", e)
		}
		service.wg.Add(1)
		go http.Serve(l, nil)
	} else {
		err = errors.New("Server bus already started")
	}
	return err
}

func (networkBus *NetworkBus) Stop() {
	service := networkBus.service
	if service.started {
		service.wg.Done()
		service.started = false
	}
}
