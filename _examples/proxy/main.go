package main

import (
	"net/http"
	"path"

	"github.com/apex/go-apex"
	"github.com/apex/go-apex/proxy"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/hello/", handleHello)

	apex.Handle(proxy.Serve(mux))
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Root"))
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	_, name := path.Split(r.URL.Path)
	w.Write([]byte("Hello " + name))
}
