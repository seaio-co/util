package endless

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

const (
	PRE_SIGNAL = iota
	POST_SIGNAL

	STATE_INIT
	STATE_RUNNING
	STATE_SHUTTING_DOWN
	STATE_TERMINATE
)

var (
	runningServerReg     sync.RWMutex
	runningServers       map[string]*endlessServer
	runningServersOrder  []string
	socketPtrOffsetMap   map[string]uint
	runningServersForked bool

	DefaultReadTimeOut    time.Duration
	DefaultWriteTimeOut   time.Duration
	DefaultMaxHeaderBytes int
	DefaultHammerTime     time.Duration

	isChild     bool
	socketOrder string

	hookableSignals []os.Signal
)

func init() {
	runningServerReg = sync.RWMutex{}
	runningServers = make(map[string]*endlessServer)
	runningServersOrder = []string{}
	socketPtrOffsetMap = make(map[string]uint)

	DefaultMaxHeaderBytes = 0 // use http.DefaultMaxHeaderBytes - which currently is 1 << 20 (1MB)

	// after a restart the parent will finish ongoing requests before
	// shutting down. set to a negative value to disable
	DefaultHammerTime = 60 * time.Second

	hookableSignals = []os.Signal{
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	}
}

type endlessServer struct {
	http.Server
	EndlessListener  net.Listener
	SignalHooks      map[int]map[os.Signal][]func()
	tlsInnerListener *endlessListener
	wg               sync.WaitGroup
	sigChan          chan os.Signal
	isChild          bool
	state            uint8
	lock             *sync.RWMutex
	BeforeBegin      func(add string)
}

/*
NewServer returns an intialized endlessServer Object. Calling Serve on it will
actually "start" the server.
*/
func NewServer(addr string, handler http.Handler) (srv *endlessServer) {
	runningServerReg.Lock()
	defer runningServerReg.Unlock()

	socketOrder = os.Getenv("ENDLESS_SOCKET_ORDER")
	isChild = os.Getenv("ENDLESS_CONTINUE") != ""

	if len(socketOrder) > 0 {
		for i, addr := range strings.Split(socketOrder, ",") {
			socketPtrOffsetMap[addr] = uint(i)
		}
	} else {
		socketPtrOffsetMap[addr] = uint(len(runningServersOrder))
	}

	srv = &endlessServer{
		wg:      sync.WaitGroup{},
		sigChan: make(chan os.Signal),
		isChild: isChild,
		SignalHooks: map[int]map[os.Signal][]func(){
			PRE_SIGNAL: map[os.Signal][]func(){
				syscall.SIGHUP:  []func(){},
				syscall.SIGUSR1: []func(){},
				syscall.SIGUSR2: []func(){},
				syscall.SIGINT:  []func(){},
				syscall.SIGTERM: []func(){},
				syscall.SIGTSTP: []func(){},
			},
			POST_SIGNAL: map[os.Signal][]func(){
				syscall.SIGHUP:  []func(){},
				syscall.SIGUSR1: []func(){},
				syscall.SIGUSR2: []func(){},
				syscall.SIGINT:  []func(){},
				syscall.SIGTERM: []func(){},
				syscall.SIGTSTP: []func(){},
			},
		},
		state: STATE_INIT,
		lock:  &sync.RWMutex{},
	}

	srv.Server.Addr = addr
	srv.Server.ReadTimeout = DefaultReadTimeOut
	srv.Server.WriteTimeout = DefaultWriteTimeOut
	srv.Server.MaxHeaderBytes = DefaultMaxHeaderBytes
	srv.Server.Handler = handler

	srv.BeforeBegin = func(addr string) {
		log.Println(syscall.Getpid(), addr)
	}

	runningServersOrder = append(runningServersOrder, addr)
	runningServers[addr] = srv

	return
}

/*
ListenAndServe listens on the TCP network address addr and then calls Serve
with handler to handle requests on incoming connections. Handler is typically
nil, in which case the DefaultServeMux is used.
*/
func ListenAndServe(addr string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServe()
}

/*
ListenAndServeTLS acts identically to ListenAndServe, except that it expects
HTTPS connections. Additionally, files containing a certificate and matching
private key for the server must be provided. If the certificate is signed by a
certificate authority, the certFile should be the concatenation of the server's
certificate followed by the CA's certificate.
*/
func ListenAndServeTLS(addr string, certFile string, keyFile string, handler http.Handler) error {
	server := NewServer(addr, handler)
	return server.ListenAndServeTLS(certFile, keyFile)
}

func (srv *endlessServer) getState() uint8 {
	srv.lock.RLock()
	defer srv.lock.RUnlock()

	return srv.state
}

func (srv *endlessServer) setState(st uint8) {
	srv.lock.Lock()
	defer srv.lock.Unlock()

	srv.state = st
}

/*
Serve accepts incoming HTTP connections on the listener l, creating a new
service goroutine for each. The service goroutines read requests and then call
handler to reply to them. Handler is typically nil, in which case the
DefaultServeMux is used.
In addition to the stl Serve behaviour each connection is added to a
sync.Waitgroup so that all outstanding connections can be served before shutting
down the server.
*/
func (srv *endlessServer) Serve() (err error) {
	defer log.Println(syscall.Getpid(), "Serve() returning...")
	srv.setState(STATE_RUNNING)
	err = srv.Server.Serve(srv.EndlessListener)
	log.Println(syscall.Getpid(), "Waiting for connections to finish...")
	srv.wg.Wait()
	srv.setState(STATE_TERMINATE)
	return
}

/*
ListenAndServe listens on the TCP network address srv.Addr and then calls Serve
to handle requests on incoming connections. If srv.Addr is blank, ":http" is
used.
*/
func (srv *endlessServer) ListenAndServe() (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	go srv.handleSignals()

	l, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return
	}

	srv.EndlessListener = newEndlessListener(l, srv)

	if srv.isChild {
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM)
	}

	srv.BeforeBegin(srv.Addr)

	return srv.Serve()
}

/*
ListenAndServeTLS listens on the TCP network address srv.Addr and then calls
Serve to handle requests on incoming TLS connections.
Filenames containing a certificate and matching private key for the server must
be provided. If the certificate is signed by a certificate authority, the
certFile should be the concatenation of the server's certificate followed by the
CA's certificate.
If srv.Addr is blank, ":https" is used.
*/
func (srv *endlessServer) ListenAndServeTLS(certFile, keyFile string) (err error) {
	addr := srv.Addr
	if addr == "" {
		addr = ":https"
	}

	config := &tls.Config{}
	if srv.TLSConfig != nil {
		*config = *srv.TLSConfig
	}
	if config.NextProtos == nil {
		config.NextProtos = []string{"http/1.1"}
	}

	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	go srv.handleSignals()

	l, err := srv.getListener(addr)
	if err != nil {
		log.Println(err)
		return
	}

	srv.tlsInnerListener = newEndlessListener(l, srv)
	srv.EndlessListener = tls.NewListener(srv.tlsInnerListener, config)

	if srv.isChild {
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM)
	}

	log.Println(syscall.Getpid(), srv.Addr)
	return srv.Serve()
}

/*
getListener either opens a new socket to listen on, or takes the acceptor socket
it got passed when restarted.
*/
func (srv *endlessServer) getListener(laddr string) (l net.Listener, err error) {
	if srv.isChild {
		var ptrOffset uint = 0
		runningServerReg.RLock()
		defer runningServerReg.RUnlock()
		if len(socketPtrOffsetMap) > 0 {
			ptrOffset = socketPtrOffsetMap[laddr]
			// log.Println("laddr", laddr, "ptr offset", socketPtrOffsetMap[laddr])
		}

		f := os.NewFile(uintptr(3+ptrOffset), "")
		l, err = net.FileListener(f)
		if err != nil {
			err = fmt.Errorf("net.FileListener error: %v", err)
			return
		}
	} else {
		l, err = net.Listen("tcp", laddr)
		if err != nil {
			err = fmt.Errorf("net.Listen error: %v", err)
			return
		}
	}
	return
}

/*
handleSignals listens for os Signals and calls any hooked in function that the
user had registered with the signal.
*/
func (srv *endlessServer) handleSignals() {
	var sig os.Signal

	signal.Notify(
		srv.sigChan,
		hookableSignals...,
	)

	pid := syscall.Getpid()
	for {
		sig = <-srv.sigChan
		srv.signalHooks(PRE_SIGNAL, sig)
		switch sig {
		case syscall.SIGHUP:
			log.Println(pid, "Received SIGHUP. forking.")
			err := srv.fork()
			if err != nil {
				log.Println("Fork err:", err)
			}
		case syscall.SIGUSR1:
			log.Println(pid, "Received SIGUSR1.")
		case syscall.SIGUSR2:
			log.Println(pid, "Received SIGUSR2.")
			srv.hammerTime(0 * time.Second)
		case syscall.SIGINT:
			log.Println(pid, "Received SIGINT.")
			srv.shutdown()
		case syscall.SIGTERM:
			log.Println(pid, "Received SIGTERM.")
			srv.shutdown()
		case syscall.SIGTSTP:
			log.Println(pid, "Received SIGTSTP.")
		default:
			log.Printf("Received %v: nothing i care about...\n", sig)
		}
		srv.signalHooks(POST_SIGNAL, sig)
	}
}

func (srv *endlessServer) signalHooks(ppFlag int, sig os.Signal) {
	if _, notSet := srv.SignalHooks[ppFlag][sig]; !notSet {
		return
	}
	for _, f := range srv.SignalHooks[ppFlag][sig] {
		f()
	}
	return
}
