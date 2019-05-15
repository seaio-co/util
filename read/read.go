package read

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Read 读取文件信息
func Read(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		fmt.Println("read_err=", err)
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return string(fd)
}
