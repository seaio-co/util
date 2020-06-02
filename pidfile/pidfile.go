package pidfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// PIDFile is a file used to store the process ID of a running process.
type PIDFile struct {
	path string
}

func checkPIDFileAlreadyExists(path string) error {
	if pidByte, err := ioutil.ReadFile(path); err == nil {
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil {
			if processExists(pid) {
				return fmt.Errorf("pid [%s] file found: %s", string(pidByte), path)
			}
		}

	}
	return nil
}

// New creates a PIDfile using the specified path.
func New(path string) (*PIDFile, error) {
	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return nil, err
	}

	return &PIDFile{path: path}, nil
}

// Remove removes the PIDFile.
func (file PIDFile) Remove() error {
	return os.Remove(file.path);
}

// ReadPIDFile 读取pid文件
func ReadPIDFile(pidfile string) (pid int, err error) {
	f, err := os.Open(pidfile)
	if err != nil {
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	pid64, err := strconv.ParseInt(string(b), 10, 64)
	pid = int(pid64)
	return
}
