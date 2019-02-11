package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func main() {

	// Routing of end points to functions (handlers)
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
	found, resetToZero := getVar("resetToZero", req)
	if found && resetToZero == "yes" {
		mux.Lock()
		nReq = 0
		mux.Unlock()
	} else {
		mux.Lock()
		nReq++
		mux.Unlock()
	}
	// FOR JSON use: www.Header().Set("Content-Type", "application/json; charset=utf-8")
	www.Header().Set("Content-Type", "text/html; charset=utf-8")
	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, "Working.  %d requests. (Version 0.0.1)\n", nReq)
	return
}

func getVar(name string, req *http.Request) (found bool, value string) {
	method := req.Method
	if method == "POST" {
		if str := req.PostFormValue(name); str != "" {
			value = str
			found = true
		}
	} else if method == "GET" {
		if str := req.URL.Query().Get(name); str != "" {
			value = str
			found = true
		}
	}
	return
}
