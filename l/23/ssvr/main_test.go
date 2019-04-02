package main

import (
	"fmt"
	"testing"

	"github.com/pschlump/godebug"
)

func Test_BindFixer_0(t *testing.T) {

	// func BindFixer(dbType, stmt string, vars []interface{}) (modStmt string, modVars []interface{}, names []string) {
	ns, nv, _ := BindFixer(`insert into "bob" ( "b", a, "c$99" ) values ( $2, $1, 'c''c' )`, []interface{}{"aa", "bb"})

	if false {
		fmt.Printf("ns=%s nv=%s\n", ns, godebug.SVar(nv))
		// ns=insert into "bob" ( "b", a, "c$99" ) values ( ?, ?, 'c''c' ) nv=["bb","aa"]
	}

	exp := `insert into "bob" ( "b", a, "c$99" ) values ( ?, ?, 'c''c' )`
	if ns != exp {
		t.Errorf("Error invalid statment, got ->%s<-, expected ->%s<-\n", ns, exp)
	}
	exp = `["bb","aa"]`
	got := godebug.SVar(nv)
	if got != exp {
		t.Errorf("Error inv got ->%s<-, expected ->%s<-\n", got, exp)
	}
}

func Test_BindFixer_1(t *testing.T) {

	ns, nv, _ := BindFixer(`insert into "bob" ( "b", a, "c$99", "a03", "a04", "a05", "a06", "a07", "a08", "a09", "a10", "a10" ) values ( $2, $1, 'c''c', $3, $4, $5, $6, $7, $8, $9, $10, $11 )`, []interface{}{"aa", "bb", "03", "04", "05", "06", "07", "08", "09", "10", "11"})

	exp := `insert into "bob" ( "b", a, "c$99", "a03", "a04", "a05", "a06", "a07", "a08", "a09", "a10", "a10" ) values ( ?, ?, 'c''c', ?, ?, ?, ?, ?, ?, ?, ?, ? )`
	if ns != exp {
		t.Errorf("Error invalid statment, got ->%s<-, expected ->%s<-\n", ns, exp)
	}
	exp = `["bb","aa","03","04","05","06","07","08","09","10","11"]`
	got := godebug.SVar(nv)
	if got != exp {
		t.Errorf("Error invalid statment, got ->%s<-, expected ->%s<-\n", got, exp)
	}
}

func Test_BindFixer_2(t *testing.T) {

	stmt := `update "documents" set "file_name" = $2, "orig_file_name" = $3 where "id" = $1`
	ns, nv, _ := BindFixer(stmt, []interface{}{"id0", "fn", "origFn"})

	exp := `update "documents" set "file_name" = ?, "orig_file_name" = ? where "id" = ?`
	if ns != exp {
		t.Errorf("Error invalid statment, got ->%s<-, expected ->%s<-\n", ns, exp)
	}
	exp = `["fn","origFn","id0"]`
	got := godebug.SVar(nv)
	if got != exp {
		t.Errorf("Error invalid statment, got ->%s<-, expected ->%s<-\n", got, exp)
	}
}
