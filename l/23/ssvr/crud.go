package main

// This file is MIT licensed.

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pschlump/Go-FTL/server/sizlib"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/godebug"
	"github.com/pschlump/uuid"
)

type CrudConfig struct {
	URIPath             string   // Path that will reach this end point
	AuthKey             bool     // Require an auth_key
	JWTKey              bool     // Require a JWT token authentntication header (logged in)
	MethodsAllowed      []string // Set of methods that are allowed
	TableName           string   // table name
	InsertCols          []string // Valid coluns for insert
	InsertPkCol         string   // PK during insert
	UpdateCols          []string // Valid columns for update
	UpdatePkCol         string   // PK during update
	WhereCols           []string // Set of columns that can be used in the "where" clause.
	SelectRequiresWhere bool     // if true, then where must be specified -- can not return entire table.
	ProjectedCols       []string // Set of columns that are projected in a select (GET).
}
type ParamListItem struct {
	ReqVar    string // variable for GetVar()
	ParamName string // Name of variable (Info Only)
	AutoGen   bool
	Required  bool
}
type CrudStoredProcConfig struct {
	URIPath             string          // Path that will reach this end point
	AuthKey             bool            // Require an auth_key
	JWTKey              bool            // Require a JWT token authentntication header (logged in)
	StoredProcedureName string          // Name of stored procedure to call.
	TableNameList       []string        // table name update/used in call (Info Only)
	ParameterList       []ParamListItem // Pairs of values
}
type CrudQueryConfig struct {
	URIPath       string          // Path that will reach this end point
	AuthKey       bool            // Require an auth_key
	JWTKey        bool            // Require a JWT token authentntication header (logged in)
	QueryString   string          // "select ... from tables where {{.where}} {{.order_by}}
	TableNameList []string        // table name used in call (Info Only)
	ParameterList []ParamListItem // Pairs of values
}

// handleSPClosure := func(cc CrudStoredProcConfig, ii int) func(www http.ResponseWriter, req *http.Request) {
func HandleStoredProcedureConfig(www http.ResponseWriter, req *http.Request, SPData CrudStoredProcConfig, posInTable int) {

	if SPData.AuthKey {
		if !IsAuthKeyValid(www, req) {
			fmt.Printf("%sAT: %s api_key wrong\n%s", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
			return
		}
	}

	method := MethodReplace(www, req)

	switch method {

	case "GET", "POST": // select
		// select <Name> ( $1 ... $n ) as "x"
		vals, inputData, err := GetStoredProcNames(www, req, SPData.ParameterList, SPData.StoredProcedureName, SPData.URIPath)
		if db_flag["HandleCRUDSP"] {
			fmt.Printf("AT: %s vals [%s] inputData %s err %s\n", godebug.LF(), vals, godebug.SVar(inputData), err)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error setting up call to %s error %s at %s\n", SPData.StoredProcedureName, err, godebug.LF())
			return
		}
		stmt := fmt.Sprintf("select %s ( %s ) as \"x\"", SPData.StoredProcedureName, vals)
		if db_flag["HandleCRUDSP"] {
			fmt.Printf("AT: %s stmt [%s]\n", godebug.LF(), stmt)
		}
		var rawData string
		err = SqliteQueryRow(stmt, inputData...).Scan(&rawData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching return data form %s ->%s<- error %s at %s\n", SPData.StoredProcedureName, stmt, err, godebug.LF())
			www.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		if db_flag["HandleCRUDSP"] {
			fmt.Printf("AT: %s rawData ->%s<-\n", godebug.LF(), rawData)
		}
		if isTLS {
			www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		www.Header().Set("Content-Type", "application/json; charset=utf-8")
		// fmt.Fprintf(www, `{"status":"success","id":%q,"data":%s}`, id, rawData)
		fmt.Fprintf(www, rawData)
		return
	default:
		if db_flag["HandleCRUDSP"] {
			fmt.Printf("AT: %s method [%s]\n", godebug.LF(), req.Method)
		}
		www.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	www.WriteHeader(http.StatusInternalServerError) // 500
	return

}

//func MultiData ( www http.ResponseWriter, req *http.Request, CrudData CrudConfig, posInTable int) bool {
func MultiData(www http.ResponseWriter, req *http.Request) (raw string, found bool) {
	found, raw = GetVar("__data__", www, req)
	return
}

type MultiRv struct {
	Status string `json:"status"`
	Msg    string
	RowNum int
	Id     string
}

// HandleCRUDConfig uses a CrudConfig to setup a table handler for requests.
func HandleCRUDConfig(www http.ResponseWriter, req *http.Request, CrudData CrudConfig, posInTable int) {

	if db_flag["HandleCRUD"] {
		fmt.Printf("Top of HandleCrudConfig: AT: %s\n", godebug.LF())
	}

	if CrudData.AuthKey {
		if !IsAuthKeyValid(www, req) {
			fmt.Printf("%sAT: %s api_key wrong\n%s", MiscLib.ColorRed, godebug.LF(), MiscLib.ColorReset)
			return
		}
	}

	method := MethodReplace(www, req)

	if !InArray(method, CrudData.MethodsAllowed) {
		www.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	switch method {

	case "GET": // select
		var data []map[string]interface{}
		if _, found := MultiData(www, req); found {
			// concatenated set of "selects" from where clauses, could be a requst with just a bunch of "id" and pull those rows.
			// For example: { [ { "id": "123" }, { "id": "444" } ] }
			// Returns 2 rows of data.
			// xyzzy4000 -mutli- delete -- not implemented yet.
		} else {
			// if "id" - then pull just id, else select all from...
			found_id, id := GetVar("id", www, req)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s found_id %v id [%s]\n", godebug.LF(), found_id, id)
			}
			if found_id {
				stmt := fmt.Sprintf("select %s from %q where \"id\" = $1", GenProjected(CrudData.ProjectedCols), CrudData.TableName)
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt [%s]\n", godebug.LF(), stmt)
				}
				rows, err := SqliteQuery(stmt, id)
				defer rows.Close()
				data, _, _ = sizlib.RowsToInterface(rows)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching form %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
					www.WriteHeader(http.StatusInternalServerError) // 500
					return
				}
			} else if cols, colsData, found := FoundCol(www, req, CrudData.WhereCols); found {
				// maybee - page-ing should occure at this point!!
				stmt := fmt.Sprintf("select %s from %q where %s", GenProjected(CrudData.ProjectedCols), CrudData.TableName, GenWhere(cols))
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt(Where Generated) [%s] data=%s\n", godebug.LF(), stmt, godebug.SVar(colsData))
				}
				rows, err := SqliteQuery(stmt, colsData...)
				defer rows.Close()
				data, _, _ = sizlib.RowsToInterface(rows)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching form %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
					www.WriteHeader(http.StatusInternalServerError) // 500
					return
				}
				fmt.Printf("%sAT: %s%s\n", MiscLib.ColorYellow, godebug.LF(), MiscLib.ColorReset)

			} else if CrudData.SelectRequiresWhere { // If true, then were must be specified - can not do a full-table below.
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s method [%s]\n", godebug.LF(), req.Method)
				}
				www.WriteHeader(http.StatusMethodNotAllowed) // 405
				return

			} else {
				// page-ing should occure at this point!!			___page__ and page-size will be needed
				stmt := fmt.Sprintf("select * from %q order by \"created\"", CrudData.TableName)
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt [%s]\n", godebug.LF(), stmt)
				}
				rows, err := SqliteQuery(stmt)
				defer rows.Close()
				data, _, _ = sizlib.RowsToInterface(rows)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching form %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
					www.WriteHeader(http.StatusInternalServerError) // 500
					return
				}
			}
		}
		if db_flag["HandleCRUD"] {
			fmt.Printf("%sAT: %s data ->%s<-%s\n", MiscLib.ColorYellow, godebug.LF(), godebug.SVarI(data), MiscLib.ColorReset)
		}
		if isTLS {
			www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		www.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(www, `{"status":"success","data":%s}`, godebug.SVarI(data))
		return

	case "POST": // insert
		// check for __data__ - if so then do a multi-insert/update??
		if raw, found := MultiData(www, req); found {
			if db_flag["HandleCRUD"] {
				fmt.Printf("%sAT: %s raw -->>%s<<--%s\n", MiscLib.ColorYellow, godebug.LF(), raw, MiscLib.ColorReset)
			}
			MultiInsertUpdate(www, req, CrudData, raw, posInTable)
		} else {

			cols, vals, inputData, id, err := GetInsertNames(www, req, CrudData.InsertCols, CrudData.InsertPkCol)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s cols [%s] vals [%s] inputData %s id %s err %s\n", godebug.LF(), cols, vals, godebug.SVar(inputData), id, err)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting up insert %s error %s at %s\n", CrudData.TableName, err, godebug.LF())
				return
			}
			stmt := fmt.Sprintf("insert into %q ( %s ) values ( %s )", CrudData.TableName, cols, vals)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s stmt [%s] data=%s\n", godebug.LF(), stmt, godebug.SVar(vals))
			}
			err = SqliteUpdate(stmt, inputData...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error inserting to %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
				www.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			if isTLS {
				www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
			}
			www.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(www, `{"status":"success","id":%q}`, id)
		}
		return

	case "PUT": // update / insert
		// check for __data__ - if so then do a multi-insert/update??
		if raw, found := MultiData(www, req); found {
			if db_flag["HandleCRUD"] {
				fmt.Printf("%sAT: %s raw -->>%s<<--%s\n", MiscLib.ColorYellow, godebug.LF(), raw, MiscLib.ColorReset)
			}
			MultiInsertUpdate(www, req, CrudData, raw, posInTable)
		} else {
			updCols, inputData, id, err := GetUpdateNames(www, req, CrudData.UpdateCols, CrudData.UpdatePkCol)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s updCols [%s] inputData %s id %s err %s\n", godebug.LF(), updCols, godebug.SVar(inputData), id, err)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting up update %s error %s at %s\n", CrudData.TableName, err, godebug.LF()) // TODO - table-name
				return
			}
			stmt := fmt.Sprintf("update %q set %s where \"id\" = $1", CrudData.TableName, updCols)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s stmt [%s] data=%s\n", godebug.LF(), stmt, godebug.SVar(inputData))
			}
			err = SqliteUpdate(stmt, inputData...)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
				www.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			fmt.Fprintf(www, `{"status":"success","id":%q}`, id)
		}
		return

	case "DELETE": // delete
		// check for __data__ - if so then do a multi-insert/update??
		if _, found := MultiData(www, req); found {
			// xyzzy4000 -mutli- delete -- not implemented yet.
		} else {
			// Delete based on id=pk
			found_id, id := GetVar("id", www, req)
			if found_id {
				stmt := fmt.Sprintf("delete from %q where \"id\" = $1", CrudData.TableName)
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt [%s]\n", godebug.LF(), stmt)
				}
				err := SqliteUpdate(stmt, id)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting from %s stmt ->%s<- id ->%s<- error %s at %s\n", CrudData.TableName, stmt, id, err, godebug.LF())
					www.WriteHeader(http.StatusInternalServerError) // 500
					return
				}
			} else {
				www.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		if isTLS {
			www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		}
		www.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprintf(www, `{"status":"success"}`)
		return

	default:
		if db_flag["HandleCRUD"] {
			fmt.Printf("AT: %s method [%s]\n", godebug.LF(), req.Method)
		}
		www.WriteHeader(http.StatusMethodNotAllowed) // 405
		return
	}

	www.WriteHeader(http.StatusInternalServerError) // 500
	return
}

func MultiInsertUpdate(www http.ResponseWriter, req *http.Request, CrudData CrudConfig, raw string, posInTable int) {
	md := make([]map[string]string, 0, len(raw)) // data is [ { one insert/update }, { 2nd insert/update } ]
	err := json.Unmarshal([]byte(raw), &md)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing __data__=%s  for %s error %s at %s\n", md, CrudData.TableName, err, godebug.LF())
		return
	}
	rv := make([]MultiRv, 0, len(md))
	status := "success" // success, partial, error
	nErr := 0
	for ii, aMd := range md {
		if method, ok := aMd["__method__"]; (ok && method == "POST") || !ok {

			cols, vals, inputData, id, err := GetInsertNamesMulti(www, req, ii, aMd, CrudData.InsertCols, CrudData.InsertPkCol)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s cols [%s] vals [%s] inputData %s id %s err %s\n", godebug.LF(), cols, vals, godebug.SVar(inputData), id, err)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting up insert %s error %s at %s, data=%s, pos=%d\n", CrudData.TableName, err, godebug.LF(), godebug.SVar(aMd), ii)
				rv = append(rv, MultiRv{Status: "error", RowNum: ii, Msg: "Error on insert"})
				nErr++
			} else {
				stmt := fmt.Sprintf("insert into %q ( %s ) values ( %s )", CrudData.TableName, cols, vals)
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt [%s] data=%s pos=%d\n", godebug.LF(), stmt, godebug.SVar(inputData), ii)
				}
				err = SqliteUpdate(stmt, inputData...)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching form %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
					rv = append(rv, MultiRv{Status: "error", RowNum: ii, Msg: fmt.Sprintf("Error inserting to %s", CrudData.TableName)})
					nErr++
				} else {
					rv = append(rv, MultiRv{Status: "success", RowNum: ii, Id: id}) // success case, send back Id
				}
			}
		} else if ok && method == "PUT" {
			updCols, inputData, id, err := GetUpdateNamesMulti(www, req, ii, aMd, CrudData.UpdateCols, CrudData.UpdatePkCol)
			if db_flag["HandleCRUD"] {
				fmt.Printf("AT: %s updCols [%s] inputData %s id %s err %s\n", godebug.LF(), updCols, godebug.SVar(inputData), id, err)
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error setting up insert %s error %s at %s\n", CrudData.TableName, err, godebug.LF()) // TODO - table-name
				rv = append(rv, MultiRv{Status: "error", RowNum: ii, Msg: "Error on update"})
				nErr++
			} else {
				stmt := fmt.Sprintf("update %q set %s where \"id\" = $1", CrudData.TableName, updCols)
				if db_flag["HandleCRUD"] {
					fmt.Printf("AT: %s stmt [%s] data=%s pos=%d\n", godebug.LF(), stmt, godebug.SVar(inputData), ii)
				}
				err = SqliteUpdate(stmt, inputData...)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error updating %s ->%s<- error %s at %s\n", CrudData.TableName, stmt, err, godebug.LF())
					rv = append(rv, MultiRv{Status: "error", RowNum: ii, Msg: fmt.Sprintf("Error updating %s", CrudData.TableName)})
					nErr++
				} else {
					rv = append(rv, MultiRv{Status: "success", RowNum: ii, Id: id}) // success case, send back Id
				}
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error %s at %s\n", CrudData.TableName, godebug.LF())
			nErr++
			rv = append(rv, MultiRv{Status: "error", RowNum: ii, Msg: fmt.Sprintf("Invalid __method__=%s for %s, processing terminated.", method, CrudData.TableName)})
			break
		}
	}
	if nErr == 0 {
		// status = "success"
	} else if nErr < len(md) {
		status = "partial"
	} else {
		status = "error"
	}
	if isTLS {
		www.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}
	www.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(www, `{"status":%q,"data":%s}`, status, godebug.SVar(rv))
}

// cols, colsData, found := FoundCol ( www, req, CrudData.WhereCols )
func FoundCol(www http.ResponseWriter, req *http.Request, WhereCols []string) (cols []string, colsData []interface{}, found bool) {
	if len(WhereCols) == 0 {
		return
	}
	for _, col := range WhereCols {
		ok, val := GetVar(col, www, req)
		// fmt.Printf("FoundCol: col [%s] ok=%v val= ->%s<- AT: %s\n", col, ok, val, godebug.LF())
		if ok {
			// fmt.Printf("FoundCol: col [%s] AT: %s\n", col, godebug.LF())
			found = true
			cols = append(cols, col)
			colsData = append(colsData, val)
		}
	}
	return
}

// GenWhere(cols)
func GenWhere(cols []string) string {
	if len(cols) == 0 {
		return ""
	}
	com := ""
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	for dd, ss := range cols {
		fmt.Fprintf(foo, "%s%q = $%d", com, ss, dd+1)
		com = " and "
	}
	foo.Flush()
	where := b.String()
	godebug.DbPf(db_flag["HandleCRUD.GenWhere"], "%(Cyan)%(LF) where [%s]\n", where)
	return where

}

func GenProjected(ProjectedCols []string) (rv string) {
	if len(ProjectedCols) == 0 {
		return "*"
	}
	rv = ""
	com := ""
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)
	for _, ss := range ProjectedCols {
		fmt.Fprintf(foo, "%s%q", com, ss)
		com = ", "
	}
	foo.Flush()
	return b.String()
}

/*
   type ParamListItem struct {
   	ReqVar    string // variable for GetVar()
   	ParamName string // Name of variable (Info Only)
   	AutoGen   bool
	Required  bool
   }
   type CrudStoredProcConfig struct {
   	URIPath             string          // Path that will reach this end point
   	AuthKey             bool            // Require an auth_key
   	JWTKey              bool            // Require a JWT token authentntication header
   	StoredProcedureName string          // Name of stored procedure to call.
   	TableNameList       []string        // table name update/used in call (Info Only)
   	ParameterList       []ParamListItem // Pairs of values
   }
*/
// vals, inputData, id, err := GetStoredProcNames(www, req, SPData.ParameterList, SPData.StoredProcedureName, SPData.URIPath)
func GetStoredProcNames(www http.ResponseWriter, req *http.Request, potentialCols []ParamListItem, StoredProcdureName, URIPath string) (vals string, inputData []interface{}, err error) {
	inputData = make([]interface{}, 0, len(potentialCols))
	nc := 1
	com := ""
	for _, col := range potentialCols {
		colName := col.ReqVar
		found, colVal := GetVar(colName, www, req)
		if col.AutoGen && !found {
			newUUID, err1 := uuid.NewV4()
			err = err1
			if err != nil {
				err = fmt.Errorf("An error occurred generating a UUID: %s", err)
				fmt.Fprintf(os.Stderr, "Error %s", err)
				www.WriteHeader(http.StatusInternalServerError) // 500
				return
			}
			colVal = newUUID.String()
		} else if !found && col.Required {
			err = fmt.Errorf("Missing %s in call to %s - endpoint %s", colName, StoredProcdureName, URIPath)
			fmt.Fprintf(os.Stderr, "Error %s", err)
			www.WriteHeader(http.StatusInternalServerError) // 500
			return
		}

		inputData = append(inputData, colVal)
		vals += fmt.Sprintf("%s$%d", com, nc)
		nc++
		com = ", "
	}
	return
}

// cols, vals, inputData, id, err := GetInsertNamesMulti(www, req, ii, aMd, CrudData.InsertCols, CrudData.InsertPkCol)
func GetInsertNamesMulti(www http.ResponseWriter, req *http.Request, pos int, aMd map[string]string, potentialCols []string, pkCol string) (
	cols, vals string, inputData []interface{}, id string, err error,
) {
	inputData = make([]interface{}, 0, len(potentialCols))
	colsSlice := make([]string, 0, len(potentialCols))
	valsSlice := make([]string, 0, len(potentialCols))
	// found_pk, pk := GetVar(pkCol, www, req)
	pk, found_pk := aMd[pkCol]
	if !found_pk {
		newUUID, err1 := uuid.NewV4()
		err = err1
		if err != nil {
			err = fmt.Errorf("An error occurred generating a UUID: %s", err)
			fmt.Fprintf(os.Stderr, "Error %s", err)
			// www.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		pk = newUUID.String()
	}
	id = pk
	nc := 1
	inputData = append(inputData, pk)
	colsSlice = append(colsSlice, pkCol)
	valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
	nc++
	for _, colName := range potentialCols {
		if colName != pkCol {
			// found, val := GetVar(colName, www, req)
			val, found := aMd[colName]
			if found {
				inputData = append(inputData, val)
				colsSlice = append(colsSlice, colName)
				valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
				nc++
			}
		}
	}
	com := ""
	for _, aCol := range colsSlice {
		cols += fmt.Sprintf("%s%q", com, aCol)
		com = ", "
	}
	com = ""
	for _, aVal := range valsSlice {
		vals += fmt.Sprintf("%s%s", com, aVal)
		com = ", "
	}
	return
	return
}

// GetInsertNames returns the list of columns for an insert, the list of placeholders for PG substitution of values,
// the list of values, the primary key or an error.
func GetInsertNames(www http.ResponseWriter, req *http.Request, potentialCols []string, pkCol string) (cols, vals string, inputData []interface{}, id string, err error) {
	inputData = make([]interface{}, 0, len(potentialCols))
	colsSlice := make([]string, 0, len(potentialCols))
	valsSlice := make([]string, 0, len(potentialCols))
	found_pk, pk := GetVar(pkCol, www, req)
	if !found_pk {
		newUUID, err1 := uuid.NewV4()
		err = err1
		if err != nil {
			err = fmt.Errorf("An error occurred generating a UUID: %s", err)
			fmt.Fprintf(os.Stderr, "Error %s", err)
			www.WriteHeader(http.StatusInternalServerError) // 500
			return
		}
		pk = newUUID.String()
	}
	id = pk
	nc := 1
	inputData = append(inputData, pk)
	colsSlice = append(colsSlice, pkCol)
	valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
	nc++
	for _, colName := range potentialCols {
		if colName != pkCol {
			found, val := GetVar(colName, www, req)
			if found {
				inputData = append(inputData, val)
				colsSlice = append(colsSlice, colName)
				valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
				nc++
			}
		}
	}
	com := ""
	for _, aCol := range colsSlice {
		cols += fmt.Sprintf("%s%q", com, aCol)
		com = ", "
	}
	com = ""
	for _, aVal := range valsSlice {
		vals += fmt.Sprintf("%s%s", com, aVal)
		com = ", "
	}
	return
}

// updCols, inputData, id, err := GetUpdateNamesMulti(www, req, ii, aMd, CrudData.UpdateCols, CrudData.UpdatePkCol)
func GetUpdateNamesMulti(www http.ResponseWriter, req *http.Request, pos int, aMd map[string]string, potentialCols []string, pkCol string) (
	updCols string, inputData []interface{}, id string, err error,
) {

	inputData = make([]interface{}, 0, len(potentialCols))
	colsSlice := make([]string, 0, len(potentialCols))
	valsSlice := make([]string, 0, len(potentialCols))
	colNameSlice := make([]string, 0, len(potentialCols))
	// found_pk, pk := GetVar(pkCol, www, req)
	pk, found_pk := aMd[pkCol]
	if !found_pk {
		err = fmt.Errorf("PK (%s) not included in udpate", pkCol)
		fmt.Fprintf(os.Stderr, "Error %s", err)
		// www.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	id = pk
	nc := 1
	inputData = append(inputData, pk)
	colsSlice = append(colsSlice, pkCol)
	valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
	colNameSlice = append(colNameSlice, pkCol)
	nc++
	for _, colName := range potentialCols {
		if colName != pkCol {
			// found, val := GetVar(colName, www, req)
			val, found := aMd[colName]
			if found {
				colNameSlice = append(colNameSlice, colName)
				inputData = append(inputData, val)
				colsSlice = append(colsSlice, colName)
				valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
				nc++
			}
		}
	}
	// if only ID, then no update
	if nc == 1 {
		err = fmt.Errorf("No columns updated")
		fmt.Fprintf(os.Stderr, "Error %s", err)
		// www.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	com := ""
	for ii, aCol := range colsSlice {
		colName := colNameSlice[ii]
		if colName != pkCol {
			aVal := valsSlice[ii]
			updCols += fmt.Sprintf("\n\t%s%q = %s ", com, aCol, aVal)
			com = ", "
		}
	}
	return
}

// GetUpdateNmaes returns the set of update columns and the data values for running an udpate.
func GetUpdateNames(www http.ResponseWriter, req *http.Request, potentialCols []string, pkCol string) (updCols string, inputData []interface{}, id string, err error) {
	inputData = make([]interface{}, 0, len(potentialCols))
	colsSlice := make([]string, 0, len(potentialCols))
	valsSlice := make([]string, 0, len(potentialCols))
	colNameSlice := make([]string, 0, len(potentialCols))
	found_pk, pk := GetVar(pkCol, www, req)
	if !found_pk {
		err = fmt.Errorf("PK (%s) not included in udpate", pkCol)
		fmt.Fprintf(os.Stderr, "Error %s", err)
		www.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	id = pk
	nc := 1
	inputData = append(inputData, pk)
	colsSlice = append(colsSlice, pkCol)
	valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
	colNameSlice = append(colNameSlice, pkCol)
	nc++
	for _, colName := range potentialCols {
		if colName != pkCol {
			found, val := GetVar(colName, www, req)
			if found {
				colNameSlice = append(colNameSlice, colName)
				inputData = append(inputData, val)
				colsSlice = append(colsSlice, colName)
				valsSlice = append(valsSlice, fmt.Sprintf("$%d", nc))
				nc++
			}
		}
	}
	// if only ID, then no update
	if nc == 1 {
		err = fmt.Errorf("No columns updated")
		fmt.Fprintf(os.Stderr, "Error %s", err)
		www.WriteHeader(http.StatusInternalServerError) // 500
		return
	}
	com := ""
	for ii, aCol := range colsSlice {
		colName := colNameSlice[ii]
		if colName != pkCol {
			aVal := valsSlice[ii]
			updCols += fmt.Sprintf("\n\t%s%q = %s ", com, aCol, aVal)
			com = ", "
		}
	}
	return
}

// HandleTables creates a set of closure based functions for each of thie items in TableConfig.
// Each handler is for a single end point.  Look at the comments for TableConfig for how to use
// this.
func HandleTables(mux *http.ServeMux) {
	handleTableClosure := func(cc CrudConfig, ii int) func(www http.ResponseWriter, req *http.Request) {
		return func(www http.ResponseWriter, req *http.Request) {
			HandleCRUDConfig(www, req, cc, ii)
		}
	}
	for ii, cc := range TableConfig {
		fmt.Printf("End Point: %s\n", cc.URIPath)
		mux.Handle(cc.URIPath, http.HandlerFunc(handleTableClosure(cc, ii)))
	}

	// Store Procedures - for d.b.'s that have them. (Postgres, Oracle etc)
	// 	PG: select {{.StoredProcedureName}} ( $1, ... $9 ) as "x"
	handleSPClosure := func(cc CrudStoredProcConfig, ii int) func(www http.ResponseWriter, req *http.Request) {
		return func(www http.ResponseWriter, req *http.Request) {
			HandleStoredProcedureConfig(www, req, cc, ii)
		}
	}
	for ii, cc := range StoredProcConfig {
		mux.Handle(cc.URIPath, http.HandlerFunc(handleSPClosure(cc, ii)))
	}
}
