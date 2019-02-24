package main

import (
	"net/http"

	"github.com/Univ-Wyo-Education/S19-4010/a/05/bsvr/cli"
)

func getRespHandler_TemplateFunction(cc *cli.CLI) HandlerFunc {
	return func(www http.ResponseWriter, req *http.Request) {
		if !ValidateAcctPin(cc, www, req) {
			return
		}
		// xyzzy - todo
	}
}
