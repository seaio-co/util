package zerocopy

import (
	"io"
	"os"
	"syscall"

	"golang.org/x/sys/unix"
)

func (p *Pipe) bufferSize() (int, error) {
	var (
		size  uintptr
		errno syscall.Errno
	)
	err := p.wrc.Control(func(fd uintptr) {
		size, _, errno = unix.Syscall(
			unix.SYS_FCNTL,
			fd,
			unix.F_GETPIPE_SZ,
			0,
		)
	})
	if err != nil {
		return 0, err
	}
	if errno != 0 {
		return 0, os.NewSyscallError("getpipesz", errno)
	}
	return int(size), nil
}

func (p *Pipe) setBufferSize(n int) error {
	var errno syscall.Errno
	err := p.wrc.Control(func(fd uintptr) {
		_, _, errno = unix.Syscall(
			unix.SYS_FCNTL,
			fd,
			unix.F_SETPIPE_SZ,
			uintptr(n),
		)
	})
	if err != nil {
		return err
	}
	if errno != 0 {
		return os.NewSyscallError("setpipesz", errno)
	}
	return nil
}

func (p *Pipe) read(b []byte) (int, error) {
	if p.teepipe == nil {
		return p.teerd.Read(b)
	}

	var (
		copied int64
		operr  error 
		rrcerr error 
		wrcerr error 

		waitread      = false
		readready     = false
		done          = false
		writeready    = false
		waitwrite     = false
		waitreadagain = false
	)
again:
	rrcerr = p.rrc.Read(func(prfd uintptr) bool {
		wrcerr = p.teepipe.wrc.Write(func(pwfd uintptr) bool {
			copied, operr = tee(prfd, pwfd, len(b))
			if operr == unix.EAGAIN {
				if !readready {
					waitread = true
				}
				return true
			}
			done = true
			operr = os.NewSyscallError("tee", operr)
			return true
		})
		if waitread {
			readready = true
			waitread = false
			return false
		}
		return true
	})
	if rrcerr != nil || done {
		goto end
	}
	wrcerr = p.teepipe.wrc.Write(func(pwfd uintptr) bool {
		p.rrc.Read(func(prfd uintptr) bool {
			copied, operr = tee(prfd, pwfd, len(b))
			if operr == unix.EAGAIN {
				if writeready {
					waitreadagain = true
				} else {
					waitwrite = true
				}
				return true
			}
			operr = os.NewSyscallError("tee", operr)
			return true
		})
		if waitwrite {
			writeready = true
			waitwrite = false
			return false
		}
		return true
	})
	if wrcerr != nil {
		goto end
	}
	if waitreadagain {
		goto again
	}
end:
	limit := len(b)
	if copied > 0 {
		limit = int(copied)
	}
	n, err := p.teerd.Read(b[:limit])
	if wrcerr != nil {
		return n, wrcerr
	}
	if operr != nil {
		return n, operr
	}
	return n, err
}

const maxSpliceSize = 4 << 20

func (p *Pipe) readFrom(src io.Reader) (int64, error) {
	var (
		rd    io.Reader
		limit int64 = 1<<63 - 1
	)
	lr, ok := src.(*io.LimitedReader)
	if ok {
		rd = lr.R
		limit = lr.N
	} else {
		rd = src
	}
	sc, ok := rd.(syscall.Conn)
	if !ok {
		return io.Copy(p.w, src)
	}
	rrc, err := sc.SyscallConn()
	if err != nil {
		return io.Copy(p.w, src)
	}

	var (
		atEOF  bool
		moved  int64
		operr  error
		rrcerr error
		wrcerr error

		fallback      = false
		waitread      = false
		readready     = false
		writeready    = false
		waitwrite     = false
		waitreadagain = false
	)
	if lr != nil {
		defer func(v *int64) {
			lr.N -= *v
		}(&moved)
	}
again:
	ok = false
	max := maxSpliceSize
	if int64(max) > limit {
		max = int(limit)
	}
	rrcerr = rrc.Read(func(rfd uintptr) bool {
		wrcerr = p.wrc.Write(func(pwfd uintptr) bool {
			var n int
			n, operr = splice(rfd, pwfd, max)
			if n > 0 {
				limit -= int64(n)
				moved += int64(n)
			}
			if operr == unix.EINVAL {
				fallback = true
				return true
			}
			if operr == unix.EAGAIN {
				waitread = !readready
				return true
			}
			if operr == nil {
				if n == 0 {
					atEOF = true
				} else {
					ok = true
				}
			}
			operr = os.NewSyscallError("splice", operr)
			return true
		})
		if waitread {
			readready = true
			waitread = false
			return false
		}
		return true
	})
	if fallback {
		return io.Copy(p.w, src)
	}
	if wrcerr != nil || atEOF {
		return moved, wrcerr
	}
	if ok {
		if limit > 0 {
			goto again
		}
		goto end
	}
	wrcerr = p.wrc.Write(func(pwfd uintptr) bool {
		rrcerr = rrc.Read(func(rfd uintptr) bool {
			var n int
			n, operr = splice(rfd, pwfd, max)
			if n > 0 {
				limit -= int64(n)
				moved += int64(n)
			}
			if operr == unix.EAGAIN {
				if writeready {
					waitwrite = false
					waitreadagain = true
				} else {
					waitwrite = true
				}
				return true
			}
			operr = os.NewSyscallError("splice", operr)
			return true
		})
		if waitwrite {
			writeready = true
			waitwrite = false
			return false
		}
		return true
	})
	if rrcerr != nil {
		return moved, rrcerr
	}
	if wrcerr != nil {
		return moved, wrcerr
	}
	if operr != nil {
		return moved, operr
	}
	if limit > 0 || waitreadagain {
		goto again
	}
end:
	return moved, nil
}

func (p *Pipe) writeTo(dst io.Writer) (int64, error) {
	sc, ok := dst.(syscall.Conn)
	if !ok {
		return io.Copy(dst, onlyReader{p})
	}
	wrc, err := sc.SyscallConn()
	if err != nil {
		return io.Copy(dst, onlyReader{p})
	}

	var (
		atEOF  bool
		moved  int64
		operr  error
		rrcerr error
		wrcerr error

		fallback      = false
		waitread      = false
		readready     = false
		writeready    = false
		waitwrite     = false
		waitreadagain = false
	)
again:
	ok = false
	rrcerr = p.rrc.Read(func(rfd uintptr) bool {
		wrcerr = wrc.Write(func(pwfd uintptr) bool {
			var n int
			n, operr = splice(rfd, pwfd, maxSpliceSize)
			if n > 0 {
				moved += int64(n)
			}
			if operr == unix.EINVAL {
				fallback = true
				return true
			}
			if operr == unix.EAGAIN {
				if !readready {
					waitread = true
				}
				return true
			}
			if operr == nil {
				if n == 0 {
					atEOF = true
				} else {
					ok = true
				}
			}
			operr = os.NewSyscallError("splice", operr)
			return true
		})
		if waitread {
			readready = true
			waitread = false
			return false
		}
		return true
	})
	if fallback {
		return io.Copy(dst, onlyReader{p})
	}
	if wrcerr != nil || atEOF {
		return moved, wrcerr
	}
	if ok {
		goto end
	}

	wrcerr = wrc.Write(func(pwfd uintptr) bool {
		rrcerr = p.rrc.Read(func(rfd uintptr) bool {
			var n int
			n, operr = splice(rfd, pwfd, maxSpliceSize)
			if n > 0 {
				moved += int64(n)
			}
			if operr == unix.EAGAIN {
				if writeready {
					waitreadagain = true
				} else {
					waitwrite = true
				}
				return true
			}
			operr = os.NewSyscallError("splice", operr)
			return true
		})
		if waitwrite {
			writeready = true
			waitwrite = false
			return false
		}
		return true
	})
	if rrcerr != nil {
		return moved, rrcerr
	}
	if wrcerr != nil {
		return moved, wrcerr
	}
	if operr != nil {
		return moved, operr
	}
	if waitreadagain {
		goto again
	}
end:
	return moved, nil
}

func transfer(dst io.Writer, src io.Reader) (int64, error) {
	var (
		rd    io.Reader
		limit int64 = 1<<63 - 1
	)
	lr, ok := src.(*io.LimitedReader)
	if ok {
		rd = lr.R
		limit = lr.N
	} else {
		rd = src
	}
	rsc, ok := rd.(syscall.Conn)
	if !ok {
		return io.Copy(dst, src)
	}
	rrc, err := rsc.SyscallConn()
	if err != nil {
		return io.Copy(dst, src)
	}

	wsc, ok := dst.(syscall.Conn)
	if !ok {
		return io.Copy(dst, src)
	}
	wrc, err := wsc.SyscallConn()
	if err != nil {
		return io.Copy(dst, src)
	}

	p, err := NewPipe()
	if err != nil {
		return io.Copy(dst, src)
	}

	var moved int64 = 0
	if lr != nil {
		defer func(v *int64) {
			lr.N -= *v
		}(&moved)
	}
	for limit > 0 {
		max := maxSpliceSize
		if int64(max) > limit {
			max = int(limit)
		}
		inpipe, fallback, err := spliceDrain(p, rrc, max)
		limit -= int64(inpipe)
		if fallback {
			return io.Copy(dst, src)
		}
		if inpipe == 0 && err == nil {
			return moved, nil
		}
		if err != nil {
			return moved, err
		}
		n, fallback, err := splicePump(wrc, p, inpipe)
		if n > 0 {
			moved += int64(n)
		}
		if fallback {
			n1, err := io.CopyN(dst, p.w, int64(inpipe))
			moved += n1
			if err != nil {
				return n1, err
			}
			n2, err := io.Copy(dst, src)
			return n1 + n2, err
		}
		if err != nil {
			return moved, err
		}
	}
	return moved, nil
}

func spliceDrain(p *Pipe, rrc syscall.RawConn, max int) (int, bool, error) {
	var (
		moved  int
		rrcerr error
		serr   error
	)
	fallback := false
	err := p.wrc.Write(func(pwfd uintptr) bool {
		rrcerr = rrc.Read(func(rfd uintptr) bool {
			var n int
			n, serr = splice(rfd, pwfd, max)
			moved = int(n)
			if serr == unix.EINVAL {
				fallback = true
				return true
			}
			if serr == unix.EAGAIN {
				return false
			}
			serr = os.NewSyscallError("splice", serr)
			return true
		})
		return true
	})
	if err != nil {
		return 0, false, err
	}
	if rrcerr != nil {
		return 0, false, rrcerr
	}
	return moved, fallback, serr
}

func splicePump(wrc syscall.RawConn, p *Pipe, inpipe int) (int, bool, error) {
	var (
		fallback bool
		moved    int
		wrcerr   error
		serr     error
	)
again:
	err := p.rrc.Read(func(prfd uintptr) bool {
		wrcerr = wrc.Write(func(wfd uintptr) bool {
			var n int
			n, serr = splice(prfd, wfd, inpipe)
			if n > 0 {
				moved += int(n)
				inpipe -= int(n)
			}
			if serr == unix.EINVAL {
				fallback = true
				return true
			}
			if serr == unix.EAGAIN {
				return false
			}
			serr = os.NewSyscallError("splice", serr)
			return true
		})
		return true
	})
	if fallback {
		return moved, true, nil
	}
	if err != nil {
		return moved, false, err
	}
	if wrcerr != nil {
		return moved, false, wrcerr
	}
	if serr != nil {
		return moved, false, serr
	}
	if inpipe > 0 {
		goto again
	}
	return moved, false, nil
}

func (p *Pipe) tee(w io.Writer) {
	tp, ok := w.(*Pipe)
	if ok {
		p.teepipe = tp
	} else {
		p.teerd = io.TeeReader(p.r, w)
	}
}

type onlyReader struct {
	io.Reader
}

func tee(rfd, wfd uintptr, max int) (int64, error) {
	return unix.Tee(int(rfd), int(wfd), max, unix.SPLICE_F_NONBLOCK)
}

func splice(rfd, wfd uintptr, max int) (int, error) {
	n, err := unix.Splice(int(rfd), nil, int(wfd), nil, max, unix.SPLICE_F_NONBLOCK)
	return int(n), err
}
