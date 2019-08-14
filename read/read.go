package read

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
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

// ReadFile 读取文件，返回*os.File
func ReadFile(filepath string) *os.File {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(path.Dir(filepath), os.ModePerm)
		file, err = os.Create(filepath)
	}
	return file
}
