package main

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/Univ-Wyo-Education/S19-4010/a/07/eth/lib/SignedDataVersion01"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/pschlump/HashStrings"
)

// https://hanezu.github.io/posts/Enable-WebSocket-support-of-Ganache-CLI-and-Subscribe-to-Events.html

func ConnectToEthereum() (err error) {

	var client *ethclient.Client
	client, err = ethclient.Dial(gCfg.URL_WS_8546)
	if err != nil {
		return fmt.Errorf("Error connecting to Geth server: %s error:[%s]", gCfg.URL_WS_8546, err)
	}
	gCfg.Client = client

	var clientrpc *rpc.Client
	clientrpc, err = rpc.Dial(gCfg.URL_WS_8546)
	if err != nil {
		return fmt.Errorf("Error connecting to Geth server: %s error:[%s]", gCfg.URL_WS_8546, err)
	}
	gCfg.ClientRPC = clientrpc

	var clientws *rpc.Client
	clientws, err = rpc.Dial(gCfg.URL_8545)
	if err != nil {
		return fmt.Errorf("Error connecting to Geth server: %s [%s]", gCfg.URL_8545, err)
	}
	gCfg.ClientWS = clientws

	gCfg.ASignedDataContract = NewSignedData(&gCfg) // Setup Contract

	return
}

type SignedDataContract struct {
	Caller          *SignedDataVersion01.SignedDataVersion01Caller
	CallerOpts      *bind.CallOpts
	Transactor      *SignedDataVersion01.SignedDataVersion01Transactor
	TransactorOpts  *bind.TransactOpts
	Contract        *SignedDataVersion01.SignedDataVersion01
	ContractAddress common.Address
}

func NewSignedData(cfg *GlobalConfigData) (rv *SignedDataContract) {

	addrHex, ok := gCfg.ContractAddress["SignedData"]
	if !ok {
		fmt.Fprintf(os.Stderr, "Missing address for 'SignedData' in configuration (cfg.json)")
	}
	contractAddress := common.HexToAddress(addrHex)

	if gCfg.AccountKey == nil {
		key, err := DecryptKeyFile(gCfg.KeyFile, gCfg.KeyFilePassword)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read KeyFile: %s: [%s]", gCfg.KeyFile, err)
			os.Exit(1)
		}
		gCfg.AccountKey = key
	}

	caller, err := SignedDataVersion01.NewSignedDataVersion01Caller(contractAddress, gCfg.Client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up SigtnedDataCaller. error: [%v]", err)
		os.Exit(1)
	}

	callerOptions := &bind.CallOpts{
		From: contractAddress,
	}

	transactor, err := SignedDataVersion01.NewSignedDataVersion01Transactor(contractAddress, gCfg.Client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up SigtnedDataTranactor. error: [%s]", err)
		os.Exit(1)
	}

	transactorOptions := bind.NewKeyedTransactor(
		gCfg.AccountKey.PrivateKey,
		// payment data !!! xyzzy
	)

	theContract, err := SignedDataVersion01.NewSignedDataVersion01(contractAddress, gCfg.Client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Contract failed to instantiate. address: %s error: [%s]", addrHex, err)
		os.Exit(1)
	}

	return &SignedDataContract{
		Caller:          caller,
		CallerOpts:      callerOptions,
		Transactor:      transactor,
		TransactorOpts:  transactorOptions,
		ContractAddress: contractAddress,
		Contract:        theContract,
	}
}

// From Contract
//	function setData ( uint256 _app, uint256 _name, bytes32 _data ) public needMinPayment payable {
// From Go Code
// func (_SignedDataVersion01 *SignedDataVersion01Transactor) SetData(opts *bind.TransactOpts, _app *big.Int, _name *big.Int, _data [32]byte) (*types.Transaction, error) {
func (sdc *SignedDataContract) SetData(app, name, data string) (tx *types.Transaction, err error) {
	var _data [32]byte
	_app, ok := big.NewInt(0).SetString(app, 16)
	if !ok {
		return nil, fmt.Errorf("Invalid app hex value: [%s]", app)
	}
	_name, ok := big.NewInt(0).SetString(name, 16)
	if !ok {
		return nil, fmt.Errorf("Invalid name hex value: [%s]", name)
	}
	tmp := HashStrings.HashStrings(data)
	for i := 0; i < 32; i++ {
		_data[i] = tmp[i]
	}
	return sdc.Transactor.SetData(sdc.TransactorOpts, _app, _name, _data)
}

// From Contract
//	function getData ( uint256 _app, uint256 _name ) public view returns ( bytes32 ) {
// From Go Code
// func (_SignedDataVersion01 *SignedDataVersion01Caller) GetData(opts *bind.CallOpts, _app *big.Int, _name *big.Int) ([32]byte, error) {
func (sdc *SignedDataContract) GetData(app, name string) (hash []byte, err error) {
	_app, ok := big.NewInt(0).SetString(app, 16)
	if !ok {
		return nil, fmt.Errorf("Invalid app hex value: [%s]", app)
	}
	_name, ok := big.NewInt(0).SetString(name, 16)
	if !ok {
		return nil, fmt.Errorf("Invalid name hex value: [%s]", name)
	}
	hashTmp, err := sdc.Caller.GetData(sdc.CallerOpts, _app, _name)
	return hashTmp[:], err
}

// DecryptKeyFile reads in a key file decrypt it with the password.
func DecryptKeyFile(keyFile, password string) (*keystore.Key, error) {
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("Faield to read KeyFile %s [%v]", keyFile, err)
	}
	key, err := keystore.DecryptKey(data, password)
	if err != nil {
		return nil, fmt.Errorf("Decryption error %s [%v]", keyFile, err)
	}
	return key, nil
}
