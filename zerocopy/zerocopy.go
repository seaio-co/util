package zerocopy

import (
	"io"
	"os"
	"syscall"
)

type Pipe struct {
	r, w     *os.File
	rrc, wrc syscall.RawConn

	teerd   io.Reader
	teepipe *Pipe
}


func NewPipe() (*Pipe, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	rrc, err := r.SyscallConn()
	if err != nil {
		return nil, err
	}
	wrc, err := w.SyscallConn()
	if err != nil {
		return nil, err
	}
	return &Pipe{
		r:     r,
		w:     w,
		rrc:   rrc,
		wrc:   wrc,
		teerd: r,
	}, nil
}

func (p *Pipe) BufferSize() (int, error) {
	return p.bufferSize()
}

func (p *Pipe) SetBufferSize(n int) error {
	return p.setBufferSize(n)
}

func (p *Pipe) Read(b []byte) (n int, err error) {
	return p.read(b)
}

func (p *Pipe) CloseRead() error {
	return p.r.Close()
}

func (p *Pipe) Write(b []byte) (n int, err error) {
	return p.w.Write(b)
}

func (p *Pipe) CloseWrite() error {
	return p.w.Close()
}

func (p *Pipe) Close() error {
	err := p.r.Close()
	err1 := p.w.Close()
	if err != nil {
		return err
	}
	return err1
}

func (p *Pipe) ReadFrom(src io.Reader) (int64, error) {
	return p.readFrom(src)
}

func (p *Pipe) WriteTo(dst io.Writer) (int64, error) {
	return p.writeTo(dst)
}

func (p *Pipe) Tee(w io.Writer) {
	p.tee(w)
}

func Transfer(dst io.Writer, src io.Reader) (int64, error) {
	return transfer(dst, src)
}
