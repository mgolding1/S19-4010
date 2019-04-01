package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/pschlump/godebug"
	"github.com/pschlump/ms"
	template "github.com/pschlump/textTemplate"
)

// RunTemplate runs a template and get the results back as a string.
// This is the primary template runner for sending email.
func RunTemplate(TemplateFn string, name_of string, g_data map[string]interface{}) string {

	rtFuncMap := template.FuncMap{
		"Center":      ms.CenterStr,   //
		"PadR":        ms.PadOnRight,  //
		"PadL":        ms.PadOnLeft,   //
		"PicTime":     ms.PicTime,     //
		"FTime":       ms.StrFTime,    //
		"PicFloat":    ms.PicFloat,    //
		"nvl":         ms.Nvl,         //
		"Concat":      ms.Concat,      //
		"title":       strings.Title,  // The name "title" is what the function will be called in the template text.
		"ifDef":       ms.IfDef,       //
		"ifIsDef":     ms.IfIsDef,     //
		"ifIsNotNull": ms.IfIsNotNull, //
		"fmtDate":     ms.FmtDate,     //
		"fmtDateTS":   ms.FmtDateTS,   //
		"isEven":      ms.IsEven,      // func IsEven(x int) (r bool) {
	}

	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	t, err := template.New("simple-tempalte").Funcs(rtFuncMap).ParseFiles(TemplateFn)
	// t, err := template.New("simple-tempalte").ParseFiles(TemplateFn)
	if err != nil {
		fmt.Printf("Error(12004): parsing/reading template, %s, AT:%s\n", err, godebug.LF())
		return ""
	}

	err = t.ExecuteTemplate(foo, name_of, g_data)
	if err != nil {
		fmt.Fprintf(foo, "Error(12005): running template=%s, %s, AT:%s\n", name_of, err, godebug.LF())
		return ""
	}

	foo.Flush()
	s := b.String() // Fetch the data back from the buffer

	if db_flag["RunTemplate"] {
		fmt.Fprintf(os.Stdout, "Template Output is: ----->%s<----- AT: %s\n", s, godebug.LF())
	}

	return s

}
