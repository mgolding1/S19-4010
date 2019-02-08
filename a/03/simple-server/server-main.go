package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {
	http.Handle("/status", http.HandlerFunc(HandleStatus)) //
	http.Handle("/", http.FileServer(http.Dir("www")))

	log.Printf("Listening on port 3000\n")
	http.ListenAndServe(":3000", nil)
}

var nReq = 0
var mux *sync.Mutex

func init() {
	mux = &sync.Mutex{}
}

// HandleStatus - server to respond with a working message if up.
func HandleStatus(www http.ResponseWriter, req *http.Request) {
	mux.Lock()
	nReq++
	mux.Unlock()
	// FOR JSON use: www.Header().Set("Content-Type", "application/json; charset=utf-8")
	www.Header().Set("Content-Type", "text/html; charset=utf-8")
	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, "Working.  %d requests. (Version 0.0.1)\n", nReq)
	return
}
