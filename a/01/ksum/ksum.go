package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

func main() {
	flag.Parse()

	fns := flag.Args()
	if len(fns) == 0 {
		fmt.Fprintf(os.Stderr, "Usage: ./ksum [file ...]\n")
		os.Exit(1)
	}

	for _, fn := range fns {
		data, err := ioutil.ReadFile(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to read %s, error: %s\n", fn, err)
			continue
		}
		fmt.Printf("%s %x\n", fn, Keccak256(data))
	}
}

// Keccak256 calculates and returns the Keccak256 hash of the input data.
// From: https://github.com/ethereum/go-ethereum/blob/master/crypto/crypto.go lines 44...51.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
