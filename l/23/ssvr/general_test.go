package main

import (
	"fmt"
	"testing"
)

func TestKeyFileDecryption(t *testing.T) {
	keyFile := "./testdata/UTC--2018-02-15T19-57-35.216297214Z--6ffba2d0f4c8fd7961f516af43c55fe2d56f6044"
	boboFile := "./testdata/nope-nope-nope"

	tests := []struct {
		desc         string
		keyFile      string
		password     string
		errorMessage string
	}{
		{
			desc:         "correct password, file exists",
			keyFile:      keyFile,
			password:     "password",
			errorMessage: "%!s(<nil>)",
		},
		{
			desc:     "file non-existent",
			keyFile:  boboFile,
			password: "",
			errorMessage: fmt.Sprintf(
				"Faield to read KeyFile %s [open %s: no such file or directory]",
				boboFile,
				boboFile,
			),
		},
		{
			desc:     "file exists, password incorrect",
			keyFile:  keyFile,
			password: "nope-this-isn-t-it",
			errorMessage: fmt.Sprintf(
				"Decryption error %s [could not decrypt key with given passphrase]",
				keyFile,
			),
		},
	}

	for _, test := range tests {
		_, err := DecryptKeyFile(test.keyFile, test.password)
		message := fmt.Sprintf("%s", err)
		if message != test.errorMessage {
			t.Errorf("Test: %s Expected: [%s] Got: [%s]", test.desc, test.errorMessage, err)
		}
	}
}
