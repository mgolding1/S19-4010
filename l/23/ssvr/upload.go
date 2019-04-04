package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pschlump/HashStrings"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug"
)

func UploadFileClosure(pth string) func(w http.ResponseWriter, r *http.Request) {

	if !filelib.Exists(pth) {
		fmt.Printf("Missing directory [%s] creating it\n", pth)
		os.MkdirAll(pth, 0755)
	}

	// UploadFile uploads a file to the server
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Printf("AT: %s\n", godebug.LF())
			jsonResponse(w, http.StatusBadRequest, "Should be a POST request.")
			return
		}

		r.ParseMultipartForm(32 << 20)
		id := r.FormValue("id")
		fmt.Printf("Id=%s\n", id)
		file, handle, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("AT: %s err [%v]\n", godebug.LF(), err)
			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Error reading file data: %s"}`, err))
			return
		}
		defer file.Close()

		mimeType := handle.Header.Get("Content-Type")
		fmt.Printf("mimeType [%s]\n", mimeType)

		var file_name, aws_file_name, orig_file_name, file_hash, url_file_name string
		switch mimeType {
		case "image/jpeg":
			file_name, aws_file_name, orig_file_name, file_hash, url_file_name, err = saveFile(w, file, handle, ".jpg", pth)
		case "image/png":
			file_name, aws_file_name, orig_file_name, file_hash, url_file_name, err = saveFile(w, file, handle, ".png", pth)
		case "application/pdf":
			file_name, aws_file_name, orig_file_name, file_hash, url_file_name, err = saveFile(w, file, handle, ".pdf", pth)
		default:
			fmt.Printf("AT: %s\n", godebug.LF())
			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"The file format (mimeType=%s) is not valid."}`, mimeType))
			return
		}

		// push file to AWS S3
		if db_flag["push-to-aws"] {
			err = AddFileToS3(awsSession, file_name, aws_file_name)
			if err != nil {
				fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
				jsonResponse(w, http.StatusBadRequest, `{"status":"error","msg":"Failed to push file to S3."}`)
				return
			}
		}

		/*
			   update "t_paper_docs" set "file_name" = $2, "orig_file_name" = $3 where "id" = $1, id=e02587d8-8084-4856-563e-806b92c2e2e9,
					file_name=./files/e8a864cb348747641cb1be134829f754515204347874550d316520dda0b37f78.jpg, orig_file_name=alaska.jpg aws_file_name=%!s(MISSING)
			   		stmt := `update "documents" set "file_name" = $2, "orig_file_name" = $3 where "id" = $1`
			   AT: File: /Users/corwin/go/src/github.com/Univ-Wyo-Education/S19-4010/l/23/ssvr/upload.go LineNo:63 err=no such table: t_paper_docs
		*/
		stmt := `update "documents" set "file_name" = $2, "orig_file_name" = $3, "url_file_name" = $4 where "id" = $1`
		fmt.Printf("stmt = %s, id=%s, file_name=%s, orig_file_name=%s\n", stmt, id, file_name, orig_file_name)

		err = SqliteUpdate(stmt, id, file_name, orig_file_name, url_file_name)
		if err != nil {
			fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save data to PG. error=%s"}`, err))
			return
		}

		// call contract at this point.
		app := fmt.Sprintf("%x", HashStrings.HashStrings("app.signedcontract.com"))
		msgHash, signature, err := SignMessage(file_hash, gCfg.AccountKey)
		if err != nil {
			// xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO // xyzzy -- TODO
		}
		name := fmt.Sprintf("%x", msgHash)
		sig := fmt.Sprintf("%x", signature)
		if db_flag["call-contract"] {

			tx, err := gCfg.ASignedDataContract.SetData(app, name, sig)
			if err != nil {
				fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
				jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to call contract. error=%s"}`, err))
				return
			}

			txID := fmt.Sprintf("%x", tx)
			stmt = `update "documents" set "txid" = $2 where "id" = $1`
			fmt.Printf("stmt = %s, id=%s, txID=%s\n", stmt, id, txID)

			err = SqliteUpdate(stmt, id, txID)
			if err != nil {
				fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
				jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save txID to PG. error=%s"}`, err))
				return
			}

			jsonResponse(w, http.StatusCreated, fmt.Sprintf(`{"status":"success","txID":%q,"aws_file_name":%q,"id":%q}`, txID, aws_file_name, id))

		}

		/*
				txID, id, err := ledger.Create(id, fmt.Sprintf(`{"v":"2","t":"paper-reg","h":%q}`, file_hash))
				if err != nil {
					fmt.Printf("AT: %s err=%v txID=%s\n", godebug.LF(), err, txID)
					jsonResponse(w, http.StatusBadRequest, `{"status":"error","msg":"Failed to push data to Eth."}`)
					return
				}

			jsonResponse(w, http.StatusCreated, fmt.Sprintf(`{"status":"success","txID":%q,"aws_file_name":%q,"id":%q}`, txID, aws_file_name, id))
		*/

		jsonResponse(w, http.StatusCreated, fmt.Sprintf(`{"status":"success","aws_file_name":%q,"id":%q}`, aws_file_name, id))
	}
}

// 15e42502-e7a5-44e2-6920-b410b9308412
//    insert into t_paper_docs ("id" ) values ( '15e42502-e7a5-44e2-6920-b410b9308412' );

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader, ext string, pth string) (file_name, aws_file_name, orig_file_name, file_hash, url_file_name string, err error) {
	var data []byte
	data, err = ioutil.ReadAll(file)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save read file. error=%s"}`, err))
		return
	}

	orig_file_name = handle.Filename
	file_hash_byte := HashStrings.HashBytes(data)
	file_hash = fmt.Sprintf("%x", file_hash_byte)
	aws_file_name = fmt.Sprintf("%s%s", file_hash, ext)
	url_file_name = fmt.Sprintf("%s/%s%s", gCfg.URLUploadPath, file_hash, ext)
	file_name = fmt.Sprintf("%s/%s%s", pth, file_hash, ext)

	err = ioutil.WriteFile(file_name, data, 0644)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save write file. error=%s"}`, err))
		return
	}
	return
}

func HashFile(file_name string) (file_hash string, err error) {
	var data []byte
	data, err = ioutil.ReadFile(file_name)
	if err != nil {
		err = fmt.Errorf(`Failed to read file. error:%s`, err)
		return
	}

	file_hash_byte := HashStrings.HashBytes(data)
	file_hash = fmt.Sprintf("%x", file_hash_byte)
	return
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func SignMessage(message string, key *keystore.Key) (msgHash string, signature string, err error) {
	messageByte, err := hex.DecodeString(message)
	if err != nil {
		return "", "", fmt.Errorf("unabgle to decode message (invalid hex data) Error:%s", err)
	}
	rawSignature, err := crypto.Sign(signHash(messageByte), key.PrivateKey) // Sign Raw Bytes, Return hex of Raw Bytes
	if err != nil {
		return "", "", fmt.Errorf("unable to sign message Error:%s", err)
	}
	signature = hex.EncodeToString(rawSignature)
	return message, signature, nil
}

func ValidateMessage(file_name, sig string, key *keystore.Key) error {
	msg, err := HashFile(file_name)
	if db_flag["ValidateMessage"] {
		fmt.Printf("hash from reading file [%s]\n", msg)
	}
	if err != nil {
		return err
	}
	_, _, err = VerifySignature(gCfg.FromAddress, sig, msg)
	return err
}

// VerifySignature takes hex encoded addr, sig and msg and verifies that the signature matches with the address.
func VerifySignature(addr, sig, msg string) (recoveredAddress, recoveredPublicKey string, err error) {
	message, err := hex.DecodeString(msg)
	if err != nil {
		return "", "", fmt.Errorf("unabgle to decode message (invalid hex data) Error:%s", err)
	}
	if !common.IsHexAddress(addr) {
		return "", "", fmt.Errorf("invalid address: %s", addr)
	}
	address := common.HexToAddress(addr)
	signature, err := hex.DecodeString(sig)
	if err != nil {
		return "", "", fmt.Errorf("signature is not valid hex Error:%s", err)
	}

	if db_flag["ValidateMessage"] {
		fmt.Printf("AT: %s\n", godebug.LF())
	}

	recoveredPubkey, err := crypto.SigToPub(signHash([]byte(message)), signature)
	if err != nil || recoveredPubkey == nil {
		return "", "", fmt.Errorf("signature verification failed Error:%s", err)
	}
	recoveredPublicKey = hex.EncodeToString(crypto.FromECDSAPub(recoveredPubkey))
	rawRecoveredAddress := crypto.PubkeyToAddress(*recoveredPubkey)
	if db_flag["ValidateMessage"] {
		fmt.Printf("AT: %s recoveredPublicKey: [%s] recoved address [%x] compare to [%x]\n", godebug.LF(), recoveredPublicKey, rawRecoveredAddress, address)
	}
	if address != rawRecoveredAddress {
		return "", "", fmt.Errorf("signature did not verify, addresses did not match")
	}
	recoveredAddress = rawRecoveredAddress.Hex()
	return
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

func HandleValidateDocument(www http.ResponseWriter, req *http.Request) {

	if !IsAuthKeyValid(www, req) {
		fmt.Printf("%sAT: %s api_key wrong\n%s", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
		return
	}

	// parameters id= Id of the document - use to get signature info, file name etc.
	found, id := GetVar("id", www, req)
	if !found || id == "" {
		www.WriteHeader(406) // Invalid Request
		return
	}

	// xyzzy - use GetData off of chain?

	stmt := "select file_name, hash, signature from documents where id = $1"
	var file_name, msg, sig string
	err := SqliteQueryRow(stmt, id).Scan(&file_name, &msg, &sig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching return data form ->%s<- id=%s error %s at %s\n", stmt, id, err, godebug.LF())
		www.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	if db_flag["HandleValidateDocument"] {
		fmt.Printf("AT:%s document_file_name [%s] hash [%s] signature [%s]\n", godebug.LF(), file_name, msg, sig)
	}

	// Check signature, if ok then JSON {success}
	rv := `{"status":"success"}`
	err = ValidateMessage(file_name, sig, gCfg.AccountKey)
	if err != nil {
		rv = `{"status":"error","msg":"signature did not verify"}`
	}
	if db_flag["HandleValidateDocument"] {
		fmt.Printf("AT:%s rv [%s]\n", godebug.LF(), rv)
	}

	if isTLS {
		www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	www.Header().Set("Content-Type", "application/json; charset=utf-8")
	www.WriteHeader(http.StatusOK) // 200
	fmt.Fprintf(www, "%s", rv)
	return
}
