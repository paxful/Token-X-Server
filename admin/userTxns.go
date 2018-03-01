package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"

	"github.com/DaveAppleton/etherUtils"
	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/DaveAppleton/etherdb"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// curl --data '{"method":"trace_transaction","params":["0x17104ac9d3312d8c136b7f44d4b8b47852618065ebfa534bd2d3b5ef218ca1f3"],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:8545

// fail 0x45c270ebc47964cd82057493cfba7a43c98c9da41a1c4c9ef8efadd6b5baa5c1
// pass 0x896d3189e627ac9cb64b93aec3b403856656d6d8ced4486241ba25efc52a341f

// curl --data '{"method":"trace_transaction","params":["0x45c270ebc47964cd82057493cfba7a43c98c9da41a1c4c9ef8efadd6b5baa5c1"],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:8545
// curl --data '{"method":"trace_transaction","params":["0x896d3189e627ac9cb64b93aec3b403856656d6d8ced4486241ba25efc52a341f"],"id":1,"jsonrpc":"2.0"}' -H "Content-Type: application/json" -X POST localhost:8545

// PoolAccountName identifies special account
var PoolAccountName = "PAXFUL_POOL_ACCOUNT"

// PoolAccountKey is the POOL KEY - must exist
func PoolAccountKey() (key *ethKeys.AccountKey, err error) {
	key = ethKeys.NewKey("userKeys/" + PoolAccountName)
	if key.LoadKey() != nil {
		err = errors.New("User does not exist")
	}
	return
}

func keyTx(key *ethKeys.AccountKey) *bind.TransactOpts {
	return bind.NewKeyedTransactor(key.GetKey())
}

// ---- user
func userTx(userKey *ethKeys.AccountKey) *bind.TransactOpts {
	return keyTx(userKey)
}

// func oneEther() *big.Int {
// 	return new(big.Int).SetUint64(1000000000000000000)
// }

func intToBig(n int) *big.Int {
	return new(big.Int).SetUint64(uint64(n))
}

type csAction struct {
	Hash string
}

// SendEtherFromUser sends ether from a user address to either an ethereum address or a user
func SendEtherFromUser(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	destination := r.FormValue("destination")
	valueStr := r.FormValue("value")
	gasLimitStr := r.FormValue("gasLimit")
	gasPriceStr := r.FormValue("gasPrice")
	dataStr := r.FormValue("data")

	userKey := ethKeys.NewKey("userKeys/" + user)
	if userKey.LoadKey() != nil {
		adminError(w, "SendEtherFromUser", errors.New("User does not exist"))
		return
	}

	var destinationAddress common.Address
	var err error
	if strings.Compare(destination[:2], "0x") == 0 {
		destinationAddress = common.HexToAddress(destination)
	} else {
		destinationAddress, err = userAddress("userKeys/" + destination)
		if err != nil {
			adminError(w, "SendEtherFromUser", err)
			return
		}
	}
	value, ok := etherUtils.StrToEther(valueStr)
	if !ok {
		adminError(w, "SendEtherFromUser", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	gasLimit, err := strconv.Atoi(gasLimitStr)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		return
	}
	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if !ok {
		adminError(w, "SendEtherFromUser", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	data := common.Hex2Bytes(dataStr)
	tx, err := sendEthereum(userKey, destinationAddress, value, uint64(gasLimit), gasPrice, data)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		return
	}

	csA := csAction{Hash: tx.Hash().Hex()}

	fmt.Println(csA)
	err = json.NewEncoder(w).Encode(csA)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		fmt.Println(err)
	}
	return
}

// SendEtherToUser sends ether from PAXFUL to a user address
func SendEtherToUser(w http.ResponseWriter, r *http.Request) {
	destination := r.FormValue("user")
	valueStr := r.FormValue("value")
	gasLimitStr := r.FormValue("gasLimit")
	gasPriceStr := r.FormValue("gasPrice")
	dataStr := r.FormValue("data")

	paxKey := ethKeys.NewKey("adminKeys/paxful")
	if paxKey.LoadKey() != nil {
		adminError(w, "SendEtherFromUser", errors.New("User does not exist"))
		return
	}

	destinationAddress, err := userAddress("userKeys/" + destination)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		return
	}

	value, ok := etherUtils.StrToEther(valueStr)
	if !ok {
		adminError(w, "SendEtherFromUser", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	gasLimit, err := strconv.Atoi(gasLimitStr)
	if err != nil {
		adminError(w, "SendEtherToUser", err)
		return
	}
	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if !ok {
		adminError(w, "SendEtherFromUser", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	data := common.Hex2Bytes(dataStr)
	tx, err := sendEthereum(paxKey, destinationAddress, value, uint64(gasLimit), gasPrice, data)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		return
	}

	csA := csAction{Hash: tx.Hash().Hex()}

	fmt.Println(csA)
	err = json.NewEncoder(w).Encode(csA)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		fmt.Println(err)
	}
	return
}

// SendEtherFromPool - Send ether to user via pool - need a bit of extra gas ?
func SendEtherFromPool(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	destination := r.FormValue("destination")
	valueStr := r.FormValue("value")
	gasLimitStr := r.FormValue("gasLimit")
	gasPriceStr := r.FormValue("gasPrice")
	dataStr := r.FormValue("data")

	userKey := ethKeys.NewKey("userKeys/" + user)
	if userKey.LoadKey() != nil {
		adminError(w, "SendEtherFromPool", errors.New("User does not exist"))
		return
	}

	// we need to send enough gas for the second transaction as well
	// gas.limit * gas.price
	// assume that pool has enough for first transaction
	poolKey, err := PoolAccountKey()
	if err != nil {
		adminError(w, "SendEtherFromPool", err)
		return
	}

	value, ok := etherUtils.StrToEther(valueStr)
	if !ok {
		adminError(w, "SendEtherFromPool", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}

	gasLimit, err := strconv.Atoi(gasLimitStr)
	if err != nil {
		adminError(w, "SendEtherFromPool", err)
		return
	}

	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if !ok {
		adminError(w, "SendEtherFromPool", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	totalGasToSend := new(big.Int).Mul(gasPrice, big.NewInt(int64(gasLimit)))
	totalAmountToSend := new(big.Int).Add(value, totalGasToSend)
	data := common.Hex2Bytes(dataStr)
	tx, err := sendEthereum(poolKey, userKey.PublicKey(), totalAmountToSend, uint64(gasLimit), gasPrice, data)
	if err != nil {
		adminError(w, "SendEtherFromPool", err)
		return
	}
	etherdb.QueueSend(tx.Hash(), user, destination, valueStr, gasLimitStr, gasPriceStr, dataStr)
}

type txStatus struct {
	Status string
	Result bool
}

func GetTransactionStatus(w http.ResponseWriter, r *http.Request) {
	hashString := r.FormValue("hash")
	hash := common.HexToHash(hashString)

	client, err := getClient()
	if err != nil {
		clientError(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rcpt, err := client.TransactionReceipt(ctx, hash)
	if err != nil {
		clientError(w, err)
		return
	}
	txs := txStatus{Status: "OK", Result: strings.Compare(rcpt.Status, "0") != 0}
	if err = json.NewEncoder(w).Encode(txs); err != nil {
		clientError(w, err)
		return
	}
}
