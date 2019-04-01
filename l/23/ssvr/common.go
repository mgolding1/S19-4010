package main

// This file is MIT licensed.

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pschlump/godebug"
)

// if n, err = IsANumber ( page, www, req ) ; err != nil {
func IsANumber(s string, www http.ResponseWriter, req *http.Request) (nv int, err error) {
	var nn int64
	nn, err = strconv.ParseInt(s, 10, 64)
	if err != nil {
		www.WriteHeader(http.StatusBadRequest) // 400
	} else {
		nv = int(nn)
	}
	return
}

func IsAuthKeyValid(www http.ResponseWriter, req *http.Request) bool {
	if db_flag["IsAuthKeyValid"] {
		fmt.Printf("AT: %s - gCfg.AuthKey = [%s]\n", godebug.LF(), gCfg.AuthKey)
	}
	found, auth_key := GetVar("auth_key", www, req)
	if gCfg.AuthKey != "" {
		if db_flag["IsAuthKeyValid"] {
			fmt.Printf("AT: %s - configed AuthKey [%s], found=%v ?auth_key=[%s]\n", godebug.LF(), gCfg.AuthKey, found, auth_key)
		}
		if !found || auth_key != gCfg.AuthKey {
			if db_flag["IsAuthKeyValid"] {
				fmt.Printf("AT: %s\n", godebug.LF())
			}
			www.WriteHeader(http.StatusUnauthorized) // 401
			return false
		}
	}
	if db_flag["IsAuthKeyValid"] {
		fmt.Printf("AT: %s\n", godebug.LF())
	}
	return true
}

func InArray(lookFor string, inArr []string) bool {
	for _, v := range inArr {
		if lookFor == v {
			return true
		}
	}
	return false
}

func InArrayN(lookFor string, inArr []string) int {
	for i, v := range inArr {
		if lookFor == v {
			return i
		}
	}
	return -1
}

// MethodReplace returns a new method if __method__ is a get argument.  This allows for testing
// of code just using get requests.  That is very convenient from a browser.
func MethodReplace(www http.ResponseWriter, req *http.Request) (methodOut string) {
	methodOut = req.Method
	if db_flag["MethodReplace"] {
		fmt.Printf("Check __method__ AT: %s\n", godebug.LF())
	}
	found_method, method := GetVar("__method__", www, req)
	if found_method && req.Method == "GET" {
		if db_flag["MethodReplace"] {
			fmt.Printf("AT: %s method=%s\n", godebug.LF(), method)
		}
		if InArray(method, []string{"PUT", "POST", "DELETE"}) {
			if db_flag["MethodReplace"] {
				fmt.Printf("AT: %s method=%s\n", godebug.LF(), method)
			}
			return method
		}
	}
	return
}

// HandleStatus - server to respond with a working message if up.
func HandleStatus(www http.ResponseWriter, req *http.Request) {
	if isTLS {
		www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	www.Header().Set("Content-Type", "application/json; charset=utf-8")
	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, `{"status":"success"}`)
	return
}

/*
	stmt = `CREATE TABLE IF NOT EXISTS users (
		id 					INTEGER PRIMARY KEY,
		username 			TEXT,
		password 			TEXT
	)`
*/
func ValidUser(un, pw string) bool {
	// xyzzy - TODO At this point you really shoudl check v.s. the d.b.
	return true
}

func HandleLogin(www http.ResponseWriter, req *http.Request) {

	un_found, un := GetVar("username", www, req)
	pw_found, pw := GetVar("password", www, req)

	if !un_found || !pw_found {
		www.WriteHeader(406) // Invalid Request
	}

	if !ValidUser(un, pw) {
		www.WriteHeader(401) // Not Authorized
	}

	if isTLS {
		www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	www.Header().Set("Content-Type", "application/json; charset=utf-8")
	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, `{"status":"success","user_id":"1","auth_token":"1234"}`)
	return
}

// HandleExitServer - graceful server shutdown.
func HandleExitServer(www http.ResponseWriter, req *http.Request) {

	if !IsAuthKeyValid(www, req) {
		return
	}
	if isTLS {
		www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	www.Header().Set("Content-Type", "application/json; charset=utf-8")

	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, `{"status":"success"}`)

	go func() {
		// Implement graceful exit with auth_key
		fmt.Fprintf(os.Stderr, "\nShutting down the server... Received /exit-server?auth_key=...\n")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownWaitTime*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			fmt.Printf("Error on shutdown: [%s]\n", err)
		}
	}()
}
