package main

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pschlump/HashStrings"
	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug"
)

func UploadFileClosure(pth string) func(w http.ResponseWriter, r *http.Request) {

	if !filelib.Exists(pth) {
		fmt.Printf("Missing directory [%s]\n", pth)
		os.Exit(1)
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

		var file_name, aws_file_name, orig_file_name, file_hash string
		switch mimeType {
		case "image/jpeg":
			file_name, aws_file_name, orig_file_name, file_hash, err = saveFile(w, file, handle, ".jpg", pth)
		case "image/png":
			file_name, aws_file_name, orig_file_name, file_hash, err = saveFile(w, file, handle, ".png", pth)
		case "application/pdf":
			file_name, aws_file_name, orig_file_name, file_hash, err = saveFile(w, file, handle, ".pdf", pth)
		default:
			fmt.Printf("AT: %s\n", godebug.LF())
			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"The file format (mimeType=%s) is not valid."}`, mimeType))
			return
		}

		stmt := `update "t_paper_docs" set "file_name" = $2, "orig_file_name" = $3 where "id" = $1`
		fmt.Printf("stmt = %s, id=%s, file_name=%s, orig_file_name=%s aws_file_name=%s\n", stmt, id, file_name, orig_file_name)

		err = SqliteUpdate(stmt, id, file_name, orig_file_name, aws_file_name)
		if err != nil {
			fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
			jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save data to PG. error=%s"}`, err))
			return
		}

		// xyzzy - call contract at this point.
		_ = file_hash

		// itemID + Values
		/*
				txID, id, err := ledger.Create(id, fmt.Sprintf(`{"v":"2","t":"paper-reg","h":%q}`, file_hash))
				if err != nil {
					fmt.Printf("AT: %s err=%v txID=%s\n", godebug.LF(), err, txID)
					jsonResponse(w, http.StatusBadRequest, `{"status":"error","msg":"Failed to push data to Eth."}`)
					return
				}

				stmt = `update "t_paper_docs" set "txid" = $2 where "id" = $1`
				fmt.Printf("stmt = %s, id=%s, txID=%s\n", stmt, id, txID)

				err = SQLUpdate(stmt, id, txID)
				if err != nil {
					fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
					jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save txID to PG. error=%s"}`, err))
					return
				}
			jsonResponse(w, http.StatusCreated, fmt.Sprintf(`{"status":"success","txID":%q,"aws_file_name":%q,"id":%q}`, txID, aws_file_name, id))
		*/

		jsonResponse(w, http.StatusCreated, fmt.Sprintf(`{"status":"success","aws_file_name":%q,"id":%q}`, aws_file_name, id))
	}
}

// 15e42502-e7a5-44e2-6920-b410b9308412
//    insert into t_paper_docs ("id" ) values ( '15e42502-e7a5-44e2-6920-b410b9308412' );

func saveFile(w http.ResponseWriter, file multipart.File, handle *multipart.FileHeader, ext string, pth string) (file_name, aws_file_name, orig_file_name, file_hash string, err error) {
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
	file_name = fmt.Sprintf("%s/%s%s", pth, file_hash, ext)

	err = ioutil.WriteFile(file_name, data, 0644)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, fmt.Sprintf(`{"status":"error","msg":"Failed to save write file. error=%s"}`, err))
		return
	}
	return
}

func jsonResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
