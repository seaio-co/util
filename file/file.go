package file

import (
	"bytes"
	"errors"
	"github.com/seaio-co/util/stringutil"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"net/http"
	"io"
	"compress/gzip"
	"archive/tar"
	"net"
	"time"
)

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// get filepath base name
func Basename(file string) string {
	return path.Base(file)
}

// get filepath dir name
func Dir(file string) string {
	return path.Dir(file)
}

func InsureDir(path string) error {
	if IsExist(path) {
		return nil
	}
	return os.MkdirAll(path, os.ModePerm)
}

func Ext(file string) string {
	return path.Ext(file)
}

// rename file name
func Rename(file string, to string) error {
	return os.Rename(file, to)
}

// delete file
func Unlink(file string) error {
	return os.Remove(file)
}

// IsFile checks whether the path is a file,
// it returns false when it's a directory or does not exist.
func IsFile(filePath string) bool {
	f, e := os.Stat(filePath)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

// IsExist checks whether a file or directory exists.
// It returns false when the file or directory does not exist.
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullPath string, err error) {
	for _, path := range paths {
		if fullPath = filepath.Join(path, filename); IsExist(fullPath) {
			return
		}
	}
	err = errors.New(fullPath + " not found in paths")
	return
}

// get absolute filepath, based on built executable file
func RealPath(file string) (string, error) {
	if path.IsAbs(file) {
		return file, nil
	}
	wd, err := os.Getwd()
	return path.Join(wd, file), err
}

// get file modified time
func FileMTime(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.ModTime().Unix(), nil
}

// get file size as how many bytes
func FileSize(file string) (int64, error) {
	f, e := os.Stat(file)
	if e != nil {
		return 0, e
	}
	return f.Size(), nil
}

// list dirs under dirPath
func DirsUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := []string{}
	for i := 0; i < sz; i++ {
		if fs[i].IsDir() {
			name := fs[i].Name()
			if name != "." && name != ".." {
				ret = append(ret, name)
			}
		}
	}

	return ret, nil

}

// list files under dirPath
func FilesUnder(dirPath string) ([]string, error) {
	if !IsExist(dirPath) {
		return []string{}, nil
	}

	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return []string{}, err
	}

	sz := len(fs)
	if sz == 0 {
		return []string{}, nil
	}

	ret := []string{}
	for i := 0; i < sz; i++ {
		if !fs[i].IsDir() {
			ret = append(ret, fs[i].Name())
		}
	}

	return ret, nil

}

// RewriteFile rewrites the file.
func RewriteFile(filename string, fn func(content []byte) (newContent []byte, err error)) error {
	f, err := os.OpenFile(filename, os.O_RDWR, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	newContent, err := fn(content)
	if err != nil {
		return err
	}
	if bytes.Equal(content, newContent) {
		return nil
	}
	f.Seek(0, 0)
	f.Truncate(0)
	_, err = f.Write(newContent)
	return err
}

// ReplaceFile replaces the bytes selected by [start, end] with the new content.
func ReplaceFile(filename string, start, end int, newContent string) error {
	if start < 0 || (end >= 0 && start > end) {
		return nil
	}
	return RewriteFile(filename, func(content []byte) ([]byte, error) {
		if end < 0 || end > len(content) {
			end = len(content)
		}
		if start > end {
			start = end
		}
		return bytes.Replace(content, content[start:end], stringutil.StringToBytes(newContent), 1), nil
	})
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Create("./newFile")
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("upload success"))
}

// Compress 压缩 使用gzip压缩成tar.gz
func Compress(files []*os.File, dest string) (err error) {

	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, file := range files {
		name := file.Name()
		file.Close()
		fi, err := os.Stat(name)
		if err != nil {
			return err
		}

		// 信息头
		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		// 写信息头
		err = tw.WriteHeader(h)
		if err != nil {
			return err
		}
		fs, err := os.Open(name)

		if err != nil {
			return err
		}

		if _, err = io.Copy(tw, fs); err != nil {
			return err
		}
		fs.Close()
	}

	return nil
}

// CreateFile 创建文件并写入指定内容
func CreateFile(fileName, data string) (*os.File, error) {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		return nil, err
		//DockerFile文件内容写入
	}
	file.WriteString(data)

	return file, nil
}

// CreateTarFile 根据传入的文件，创建指定的tar文件
func CreateTarFile(name string, files ...*os.File) (*os.File, error) {

	tarFile := make([]*os.File, len(files))
	for i, file := range files {
		tarFile[i] = file
	}

	err := Compress(tarFile, name)
	if err != nil {
		return nil, err
	}

	err = deleteFile(files...)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(name)

	if err != nil {
		return nil, err
	}
	return file, nil
}

// deleteFile 删除文件，可传入多个
func deleteFile(files ...*os.File) error {
	for _, file := range files {
		file.Close()
		err := os.Remove(file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

// FileUpload 上传文件到指定URL
func FileUpload(URL string, file *os.File) (string, error) {
	resp, err := http.NewRequest("POST", URL, file)
	if err != nil {
		return "", err
	}

	resp.Header.Add("Content-Type", "application/tar")
	// 设置 TimeOut
	DefaultClient := http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}

	res, err := DefaultClient.Do(resp)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = os.Remove(file.Name())
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}