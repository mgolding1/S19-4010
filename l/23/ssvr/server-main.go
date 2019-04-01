package main

import (
	"context"
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/Univ-Wyo-Education/S19-4010/l/23/ssvr/ReadConfig"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug"
)

// ----------------------------------------------------------------------------------
// Notes:
//   Graceful Shutdown: https://stackoverflow.com/questions/39320025/how-to-stop-http-listenandserve
// ----------------------------------------------------------------------------------

var Cfg = flag.String("cfg", "cfg.json", "config file for this call")
var HostPort = flag.String("hostport", ":9004", "Host/Port to listen on")
var DbFlag = flag.String("db_flag", "", "Additional Debug Flags")
var TLS_crt = flag.String("tls_crt", "", "TLS Signed Publick Key")
var TLS_key = flag.String("tls_key", "", "TLS Signed Private Key")

type GlobalConfigData struct {
	TickerSeconds int `json:"ticker_seconds" default:"30"` // Time Ticker Seconds

	DBSqlite    string `json:"sqlite_name" default:"mydb.db"`
	LogFileName string `json:"log_file_name"`

	// debug flags:
	//		GetVal						Show fetching of values from GET/POST/...
	//		HandleCRUD					General CRUD
	//		HandleCRUD.GenWhere			CRUD generating "where" clause
	//		HandleCRUDSP				CRUD calling stored procedures
	//		IsAuthKeyValid				Check of 'auth_key'
	//		dump-db-flag				Dump flags turned on.
	DebugFlag string `json:"db_flag"`

	AuthKey string `json:"auth_key" default:""` // Auth key by default is turned off.

	// Defauilt file for TLS setup (Shoud include path), both must be specified.
	// These can be over ridden on the command line.
	TLS_crt string `json:"tls_crt" default:""`
	TLS_key string `json:"tls_key" default:""`

	// Path where files are temporary uploaded to
	UploadPath string `json:"upload_path" default:"./files"`

	TemplateDir string `default:"./tmpl"`

	StaticPath string `json:"static_path" default:"www"`
}

var gCfg GlobalConfigData
var logFilePtr *os.File
var logFileName = ""
var DB *sql.DB
var db_flag map[string]bool
var wg sync.WaitGroup
var httpServer *http.Server
var logger *log.Logger
var shutdownWaitTime = time.Duration(1)
var isTLS bool

func init() {
	isTLS = false
	db_flag = make(map[string]bool)
	logger = log.New(os.Stdout, "", 0)
}

func main() {

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Printf("Extra arguments are not supported [%s]\n", fns)
		os.Exit(1)
	}

	if Cfg == nil {
		fmt.Printf("--cfg is a required parameter\n")
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Read in Configuraiton
	// ------------------------------------------------------------------------------
	err := ReadConfig.ReadFile(*Cfg, &gCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read confguration: %s error %s\n", *Cfg, err)
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Logging File
	// ------------------------------------------------------------------------------
	if gCfg.LogFileName != "" {
		LogFile(gCfg.LogFileName)
	}

	// ------------------------------------------------------------------------------
	// TLS parameter check / setup
	// ------------------------------------------------------------------------------
	if *TLS_crt == "" && gCfg.TLS_crt != "" {
		TLS_crt = &gCfg.TLS_crt
	}
	if *TLS_key == "" && gCfg.TLS_key != "" {
		TLS_key = &gCfg.TLS_key
	}

	if *TLS_crt != "" && *TLS_key == "" {
		log.Fatalf("Must supply both .crt and .key for TLS to be turned on - fatal error.")
	} else if *TLS_crt == "" && *TLS_key != "" {
		log.Fatalf("Must supply both .crt and .key for TLS to be turned on - fatal error.")
	} else if *TLS_crt != "" && *TLS_key != "" {
		if !filelib.Exists(*TLS_crt) {
			log.Fatalf("Missing file ->%s<-\n", *TLS_crt)
		}
		if !filelib.Exists(*TLS_key) {
			log.Fatalf("Missing file ->%s<-\n", *TLS_key)
		}
		isTLS = true
	}

	// ------------------------------------------------------------------------------
	// Debug Flag Processing
	// ------------------------------------------------------------------------------
	if gCfg.DebugFlag != "" {
		ss := strings.Split(gCfg.DebugFlag, ",")
		// fmt.Printf("gCfg.DebugFlag ->%s<-\n", gCfg.DebugFlag)
		for _, sx := range ss {
			// fmt.Printf("Setting ->%s<-\n", sx)
			db_flag[sx] = true
		}
	}
	if *DbFlag != "" {
		ss := strings.Split(*DbFlag, ",")
		// fmt.Printf("gCfg.DebugFlag ->%s<-\n", gCfg.DebugFlag)
		for _, sx := range ss {
			// fmt.Printf("Setting ->%s<-\n", sx)
			db_flag[sx] = true
		}
	}
	if db_flag["dump-db-flag"] {
		fmt.Fprintf(os.Stderr, "%sDB Flags Enabled Are:%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)
		for x := range db_flag {
			fmt.Fprintf(os.Stderr, "%s\t%s%s\n", MiscLib.ColorGreen, x, MiscLib.ColorReset)
		}
	}

	fmt.Fprintf(os.Stderr, "%sAT %s, %s\n", MiscLib.ColorGreen, godebug.LF(), MiscLib.ColorReset)

	// ------------------------------------------------------------------------------
	// Connect to Database
	// ------------------------------------------------------------------------------
	ConnectToSqlite()
	CreateTables()

	fmt.Fprintf(os.Stderr, "%sConnected to Sqlite%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)

	// ------------------------------------------------------------------------------
	// Setup to use the ledger and send stuff to Eth.
	// ------------------------------------------------------------------------------
	SetupGeth()

	// ------------------------------------------------------------------------------
	// Setup HTTP End Points
	// ------------------------------------------------------------------------------
	mux := http.NewServeMux()
	mux.Handle("/api/v1/status", http.HandlerFunc(HandleStatus))          //
	mux.Handle("/status", http.HandlerFunc(HandleStatus))                 //
	mux.Handle("/api/v1/exit-server", http.HandlerFunc(HandleExitServer)) //
	mux.Handle("/login", http.HandlerFunc(HandleLogin))                   //	// // not a real login - just retruns success -
	mux.HandleFunc("/upload", UploadFileClosure(gCfg.UploadPath))         // URL to upload files with multi-part mime

	// temp test
	// mux.Handle("/api/v1/t_documents", http.HandlerFunc(HandleStatus)) //

	// For the list of end-points (URI Paths) see ./handle.go
	HandleTables(mux)

	mux.Handle("/", http.FileServer(http.Dir(gCfg.StaticPath)))

	//	mux.DumpPath()

	// ------------------------------------------------------------------------------
	// Setup signal capture
	// ------------------------------------------------------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// ------------------------------------------------------------------------------
	// Setup / Run the HTTP Server.
	// ------------------------------------------------------------------------------
	if isTLS {
		cfg := &tls.Config{
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		}
		// loggingHandler := NewApacheLoggingHandler(mux, os.Stderr)
		loggingHandler := NewApacheLoggingHandler(mux, logFilePtr)
		httpServer = &http.Server{
			Addr:    *HostPort,
			Handler: loggingHandler,
			// Handler:      mux,
			TLSConfig:    cfg,
			TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		}
	} else {
		// loggingHandler := NewApacheLoggingHandler(mux, os.Stderr)
		loggingHandler := NewApacheLoggingHandler(mux, logFilePtr)
		httpServer = &http.Server{
			Addr:    *HostPort,
			Handler: loggingHandler,
			// Handler: mux,
		}
	}

	go func() {
		wg.Add(1)
		defer wg.Done()
		if isTLS {
			fmt.Fprintf(os.Stderr, "%sListening on https://%s%s\n", MiscLib.ColorGreen, *HostPort, MiscLib.ColorReset)
			if err := httpServer.ListenAndServeTLS(*TLS_crt, *TLS_key); err != nil {
				logger.Fatal(err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "%sListening on http://%s%s\n", MiscLib.ColorGreen, *HostPort, MiscLib.ColorReset)
			if err := httpServer.ListenAndServe(); err != nil {
				logger.Fatal(err)
			}
		}
	}()

	// ------------------------------------------------------------------------------
	// Catch signals from [Contro-C]
	// ------------------------------------------------------------------------------
	select {
	case <-stop:
		fmt.Fprintf(os.Stderr, "\nShutting down the server... Received OS Signal...\n")
		ctx, cancel := context.WithTimeout(context.Background(), shutdownWaitTime*time.Second)
		defer cancel()
		err := httpServer.Shutdown(ctx)
		if err != nil {
			fmt.Printf("Error on shutdown: [%s]\n", err)
		}
	}

	// ------------------------------------------------------------------------------
	// Wait for HTTP server to exit.
	// ------------------------------------------------------------------------------
	wg.Wait()
}

func SetupGeth() {
	// xyzzy - some work to do.
}
