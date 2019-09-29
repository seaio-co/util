package common

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"strings"
)

// NewFormBody
func NewFormBody(values, files url.Values) (contentType string, bodyReader io.Reader, err error) {
	if len(files) == 0 {
		return "application/x-www-form-urlencoded", strings.NewReader(values.Encode()), nil
	}
	var rw = bytes.NewBuffer(make([]byte, 32*1024*len(files)))
	var bodyWriter = multipart.NewWriter(rw)
	var buf = make([]byte, 32*1024)
	var fileWriter io.Writer
	var f *os.File
	for fieldName, postfiles := range files {
		for _, fileName := range postfiles {
			fileWriter, err = bodyWriter.CreateFormFile(fieldName, fileName)
			if err != nil {
				return
			}
			f, err = os.Open(fileName)
			if err != nil {
				return
			}
			_, err = io.CopyBuffer(fileWriter, f, buf)
			f.Close()
			if err != nil {
				return
			}
		}
	}
	for k, v := range values {
		for _, vv := range v {
			bodyWriter.WriteField(k, vv)
		}
	}
	bodyWriter.Close()
	return bodyWriter.FormDataContentType(), rw, nil
}

type (
	Files map[string][]File
	File interface {
		Name() string
		Read(p []byte) (n int, err error)
	}
)

// NewFormBody2
func NewFormBody2(values url.Values, files Files) (contentType string, bodyReader io.Reader) {
	if len(files) == 0 {
		return "application/x-www-form-urlencoded", strings.NewReader(values.Encode())
	}
	var pr, pw = io.Pipe()
	var bodyWriter = multipart.NewWriter(pw)
	var fileWriter io.Writer
	var buf = make([]byte, 32*1024)
	go func() {
		for fieldName, postfiles := range files {
			for _, file := range postfiles {
				fileWriter, _ = bodyWriter.CreateFormFile(fieldName, file.Name())
				io.CopyBuffer(fileWriter, file, buf)
			}
		}
		for k, v := range values {
			for _, vv := range v {
				bodyWriter.WriteField(k, vv)
			}
		}
		bodyWriter.Close()
		pw.Close()
	}()
	return bodyWriter.FormDataContentType(), pr
}

// NewFile creates a file for HTTP form.
func NewFile(name string, bodyReader io.Reader) File {
	return &fileReader{name, bodyReader}
}

// fileReader file name and bytes.
type fileReader struct {
	name       string
	bodyReader io.Reader
}

func (f *fileReader) Name() string {
	return f.name
}

func (f *fileReader) Read(p []byte) (int, error) {
	return f.bodyReader.Read(p)
}

// NewJSONBody returns JSON request content type and body reader.
func NewJSONBody(v interface{}) (contentType string, bodyReader io.Reader, err error) {
	b, err := json.Marshal(v)
	if err != nil {
		return
	}
	return "application/json;charset=utf-8", bytes.NewReader(b), nil
}

// NewXMLBody returns XML request content type and body reader.
func NewXMLBody(v interface{}) (contentType string, bodyReader io.Reader, err error) {
	b, err := xml.Marshal(v)
	if err != nil {
		return
	}
	return "application/xml;charset=utf-8", bytes.NewReader(b), nil
}
