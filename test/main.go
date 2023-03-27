package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/spkg/zipfs"
)

func load(name string) (*zipfs.FileSystem, error) {
	zf, err := os.Open(fmt.Sprintf("../testdata/%s.zip", name))
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

func wrapResult(id interface{}, result interface{}, err error) interface{} {
	ret := map[string]interface{}{
		"id":      id,
		"jsonrpc": "2.0",
	}
	if err != nil {
		ret["error"] = err.Error()
	} else {
		ret["result"] = result
	}
	return ret
}

func writeResult(w http.ResponseWriter, id interface{}, result interface{}, err error) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(wrapResult(id, result, err)); err != nil {
		log.Println(err)
	}
}

func rpcService(w http.ResponseWriter, r *http.Request) {
	m := &map[string]interface{}{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(m); err != nil {
		writeResult(w, 0, nil, err)
	} else {
		params := ToArray(m, "params")
		id := AsUint64(m, "id")
		method := ToString(m, "method")

		log.Println("rpc", id, r.RequestURI, method, params)
		result, err := process(method, params)
		if err != nil {
			writeResult(w, id, nil, err)
		} else {
			writeResult(w, id, result, nil)
		}
		log.Println("response", id, err)
	}
}

func AsUint64(m *map[string]interface{}, key string) uint64 {
	item, ok := (*m)[key]
	if ok {
		i, ok := item.(float64)
		if ok {
			return uint64(i)
		}
	}
	return 0
}

func ToString(m *map[string]interface{}, key string) string {
	item, ok := (*m)[key]
	if ok {
		s, ok := item.(string)
		if ok {
			return s
		}
	}
	return ""
}

func ToArray(m *map[string]interface{}, key string) []interface{} {
	if m != nil {
		data := (*m)[key]
		if data != nil {
			list, ok := data.([]interface{})
			if ok {
				return list
			}
		}
	}
	return nil
}

func process(method string, params []interface{}) (interface{}, error) {
	return params, nil
}

func main() {
	rpcAddress := "localhost"
	rpcPort := 8080
	address := fmt.Sprintf("%s:%d", rpcAddress, rpcPort)

	http.HandleFunc("/v1/jsonrpc", rpcService)
	http.HandleFunc("/", zipfs.FileServerWith(load).ServeHTTP)

	log.Println("rpc server started on", address)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
