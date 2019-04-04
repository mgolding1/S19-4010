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
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pschlump/HashStrings"
	"github.com/pschlump/MiscLib"
	"github.com/pschlump/filelib"
	"github.com/pschlump/godebug"
	// "github.com/Univ-Wyo-Education/S19-4010/a/07/eth/lib/SignedDataVersion01"
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

	// S3 login/config options.
	// Also uses
	//		AWS_ACCESS_KEY_ID=AKIAJZDMAULPRXT7CVWA		((example))
	//		AWS_SECRET_KEY=........
	S3_bucket string `json:"s3_bucket" default:"acb-document"`
	S3_region string `json:"s3_bucket" default:"$ENV$AWS_REGION"`

	// Defauilt file for TLS setup (Shoud include path), both must be specified.
	// These can be over ridden on the command line.
	TLS_crt string `json:"tls_crt" default:""`
	TLS_key string `json:"tls_key" default:""`

	// Path where files are temporary uploaded to
	UploadPath    string `json:"upload_path" default:"./www/files"`
	URLUploadPath string `json:"url_upload_path" default:"/files"`

	TemplateDir string `default:"./tmpl"`

	StaticPath string `json:"static_path" default:"www"`

	URL_WS_8546     string            `json:"geth_ws_8546" default:"ws://127.0.0.1:9545"`         // example: ws://192.168.0.200:8546
	URL_8545        string            `json:"geth_rpc_8545" default:"http://127.0.0.1:9545"`      // example: "http://192.168.0.200:8545".
	ContractAddress map[string]string `json:"ContractAddress"`                                    // Contract names to contract addresses map
	FromAddress     string            `json:"FromAddress"`                                        // Address of account to pull funds from - this is the signing account
	KeyFilePassword string            `json:"key_file_password" default:"$ENV$Key_File_Password"` // Password to access KeyFile
	KeyFile         string            `json:"key_file" default:"$ENV$Key_File"`                   // File name for pub/priv key for Address

	Client    *ethclient.Client `json:"-"` // used in secalling contract
	ClientRPC *rpc.Client       `json:"-"`
	ClientWS  *rpc.Client       `json:"-"`

	AccountKey *keystore.Key `json:"-"`

	ASignedDataContract *SignedDataContract `json:"-"`

	BaseURLTmpl string `json:"base_url_tmpl" default:"http://127.0.0.1/base.html?hash={{.hash}}"`
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
var awsSession *session.Session

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
	// Setup connection to S3
	// ------------------------------------------------------------------------------
	if db_flag["push-to-aws"] {
		awsSession = SetupS3(gCfg.S3_bucket, gCfg.S3_region)
	}

	// ------------------------------------------------------------------------------
	// Setup to use the ledger and send stuff to Eth.
	// ------------------------------------------------------------------------------
	err = ConnectToEthereum()
	if err != nil {
		fmt.Printf("Error: %s on connecting to Geth/Ethereum - fatal.\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "%sConnected to Ethereum (Geth or Ganache)%s\n", MiscLib.ColorGreen, MiscLib.ColorReset)

	// ------------------------------------------------------------------------------
	// test code for contract calls - to verify we have contrct working.
	// ------------------------------------------------------------------------------
	if db_flag["test-contrct-calls"] {

		if db_flag["test-set-data"] {
			app := fmt.Sprintf("%x", HashStrings.HashStrings("app.signedcontract.com"))
			name := "2b2ad1becfddd20ccdefaab5e9fd160512d29cf3"
			sig := "28e7f16cc55c6b2f553e28830e62eb693ba630a1b951f5cf65dd61b9f2fb8b19db95e37afcf5eab20"
			tx, err := gCfg.ASignedDataContract.SetData(app, name, sig)
			if err != nil {
				fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
				return
			}

			txID := fmt.Sprintf("%x", tx)
			fmt.Printf("txID [%s] tx [%s] err %s\n", txID, tx, err)
		}

		if db_flag["test-get-data"] {
			app := fmt.Sprintf("%x", HashStrings.HashStrings("app.signedcontract.com"))
			name := "2b2ad1becfddd20ccdefaab5e9fd160512d29cf3"
			expected_sig := "28e7f16cc55c6b2f553e28830e62eb693ba630a1b951f5cf65dd61b9f2fb8b19db95e37afcf5eab20"
			sig, err := gCfg.ASignedDataContract.GetData(app, name)
			if err != nil {
				fmt.Printf("AT: %s err=%v\n", godebug.LF(), err)
				return
			}

			match := fmt.Sprintf("%sno%s", MiscLib.ColorRed, MiscLib.ColorReset)
			if sig == expected_sig {
				match = fmt.Sprintf("%syes%s", MiscLib.ColorGreen, MiscLib.ColorReset)
			}
			fmt.Printf("sig [%s] expected [%s] matched %s err %s\n", sig, expected_sig, match, err)
		}

		os.Exit(0)
	}

	// ------------------------------------------------------------------------------
	// Setup HTTP End Points
	// ------------------------------------------------------------------------------
	mux := http.NewServeMux()
	mux.Handle("/api/v1/status", http.HandlerFunc(HandleStatus))          //
	mux.Handle("/status", http.HandlerFunc(HandleStatus))                 //
	mux.Handle("/api/v1/exit-server", http.HandlerFunc(HandleExitServer)) //
	mux.Handle("/login", http.HandlerFunc(HandleLogin))                   // !! not a real login - just retruns success !!
	mux.HandleFunc("/upload", UploadFileClosure(gCfg.UploadPath))         // URL to upload files with multi-part mime data type

	mux.Handle("/api/v1/save-data", http.HandlerFunc(HandleStatus))   // xyzzy - TODO - save data to on-chain (uses setData) new app+name+hash
	mux.Handle("/api/v1/add-to-data", http.HandlerFunc(HandleStatus)) // xyzzy - TODO - Add a new document to chain with "EventName" Some events are final and can not be repeated.
	mux.Handle("/api/v1/get-data", http.HandlerFunc(HandleStatus))    // xyzzy - TODO - get data given app and name (uses getData) - app+name -> hash
	mux.Handle("/api/v1/get-url", http.HandlerFunc(HandleStatus))     // xyzzy - TODO - return URL given app and Name (uses getData)
	// Uses gCfg.BaseURLTmpl string `json:"base_url_tmpl" default:"http://127.0.0.1/base.html?hash={{.hash}}"`

	// For the list of end-points (URI Paths) see ./handle.go
	HandleTables(mux)

	mux.Handle("/", http.FileServer(http.Dir(gCfg.StaticPath)))

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
