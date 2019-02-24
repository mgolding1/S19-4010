package main

// Copyright 2017 The go-ethereum Authors
// Copyright 2018 Philip Schlump.
// This file uses code that was part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// The LICENSE.main_gpl3 file in this directory applies to this code.

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pborman/uuid"
)

// GenerateSignature uses a keyfile and password to sign a message.  If the input message is "" then a random message
// will be generated.  The messgae and the signature are returned.
func GenerateSignature(keyFile, password string, inMessage []byte) (message, signature string, err error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return "", "", fmt.Errorf("unable to read keyfile %s Error:%s", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return "", "", fmt.Errorf("unable to decrypt %s Error:%s", keyFile, err)
	}
	if len(inMessage) == 0 {
		inMessage, err = GenRandBytes(20)
		if err != nil {
			return "", "", fmt.Errorf("unable to generate random message Error:%s", err)
		}
	}
	message = hex.EncodeToString(inMessage)
	rawSignature, err := crypto.Sign(signHash(inMessage), key.PrivateKey) // Sign Raw Bytes, Return hex of Raw Bytes
	if err != nil {
		return "", "", fmt.Errorf("unable to sign message Error:%s", err)
	}
	signature = hex.EncodeToString(rawSignature)
	return message, signature, nil
}

// GenerateNewKeyFile will generate a new key file.
func GenerateNewKeyFile(passphrase string) error {
	var privateKey *ecdsa.PrivateKey
	var err error

	// Generate random key.
	privateKey, err = crypto.GenerateKey()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to generate random private key: %v", err)
		os.Exit(1)
	}

	// Create the keyfile object with a random UUID.
	id := uuid.NewRandom()
	key := &keystore.Key{
		Id:         id,
		Address:    crypto.PubkeyToAddress(privateKey.PublicKey),
		PrivateKey: privateKey,
	}

	// Encrypt key with passphrase.
	keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encrypting key: %v", err)
		os.Exit(1)
	}

	address := key.Address.Hex()

	newname := fmt.Sprintf("UTC--%s--%s", time.Now().UTC().Format("2006-01-02T15-04-05.9999999999Z"), address[2:])

	path := filepath.Join(GCfg.WalletPath, newname)

	// check if file already exists - if so then cowerdly refuse to ovewrite it.
	if Exists(path) {
		fmt.Fprintf(os.Stderr, "File [%s] already exists - will not overwrite\n", path)
		os.Exit(1)
	}

	// Output the file.
	if err := ioutil.WriteFile(path, keyjson, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write keyfile to %s: %v", path, err)
		os.Exit(1)
	}

	fmt.Println("Address:", address)
	fmt.Println("File Name:", path)
	return nil
}

// signHash is a helper function that calculates a hash for the given message
// that can be safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func getPassphrase(passwordIn string) (password string) {
	if passwordIn != "" {
		password = passwordIn
		return
	}

	password, err := console.Stdin.PromptPassword("Password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read password: %s", err)
		os.Exit(1)
	}
	password = strings.TrimRight(password, " \r\n\f\t")
	password = strings.TrimLeft(password, " \r\n\f\t")
	return
}
