//Example server using the ServeMux HTTP request multiplexer
//Code from Ilya Glotov's blog is used here -> https://ilyaglotov.com/blog/servemux-and-path-traversal

package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

const root = "/tmp"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filename := filepath.Join(root, strings.Trim(r.URL.Path, "/"))
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(contents)
	})

	server := &http.Server{
		Addr:    "127.0.0.1:50000",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
