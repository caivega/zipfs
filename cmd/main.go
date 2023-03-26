package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spkg/zipfs"
)

func load(name string) (*zipfs.FileSystem, error) {
	if name == "testdata" {
		zf, err := os.Open("../testdata/testdata.zip")
		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(zf)
		if err != nil {
			panic(err)
		}
		reader := bytes.NewReader(data)
		fs, err := zipfs.NewFromReaderAt(reader, int64(len(data)), nil)
		if err != nil {
			panic(err)
		}
		return fs, nil
	}
	return nil, errors.New("error")
}

func main() {
	http.HandleFunc("/", zipfs.FileServerWith(load).ServeHTTP)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
