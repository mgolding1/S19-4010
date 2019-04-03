#!/bin/bash

# check-json-syntax can be skipped - or - pull from:
# 	$ git clone https://github.com/pschlump/check-json-syntax.git
# then compile it with Go.

curl -X POST --data '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}' http://127.0.0.1:9545 | \
	check-json-syntax -p 


