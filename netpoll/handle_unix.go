package netpoll

import "syscall"

func setNonblock(fd int, nonblocking bool) (err error) {
	return syscall.SetNonblock(fd, nonblocking)
}
