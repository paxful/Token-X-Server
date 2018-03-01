package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"context"

	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/DaveAppleton/parityclient"
	"github.com/spf13/viper"
)

var client *parityclient.Client

func getClient() (client *parityclient.Client, err error) {
	if client != nil {
		return client, nil
	}
	hostChain := viper.GetString("HOST")
	endPoint := viper.GetString(hostChain)

	fmt.Println("Using ", hostChain)

	if len(endPoint) == 0 {
		endPoint = "/Users/daveappleton/Library/Ethereum/geth.ipc"
	}
	client, err = parityclient.Dial(endPoint)
	return
}

type clErr struct {
	Status string
	Error  string
}

func clientError(w http.ResponseWriter, err error) {
	cl := clErr{"ERROR", err.Error()}
	json.NewEncoder(w).Encode(cl)
}

func getUserEtherBalance(user string) (bal *big.Int, err error) {

	client, err = getClient()
	if err != nil {
		//clientError(w, err)
		return
	}

	key := ethKeys.NewKey(user)
	if key.LoadKey() != nil {
		err = errors.New("User does not exist : " + user)
		//clientError(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	fmt.Println(key.PublicKeyAsHexString())
	bal, err = client.BalanceAt(ctx, key.PublicKey(), nil)
	return
}

// SyncProgress tells you whether teh blockchain is up to date
func SyncProgress(w http.ResponseWriter, r *http.Request) {
	client, err := getClient()
	if err != nil {
		clientError(w, err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	sp, err := client.SyncProgress(ctx)
	if err != nil {
		clientError(w, err)
		return
	}
	json.NewEncoder(w).Encode(sp)
}

// CheckUserEtherBalance - get users eth balance
func CheckUserEtherBalance(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if !validateUser(w, user) {
		userError(w, user, "", fmt.Errorf("invalid user"))
		return
	}
	bal, err := getUserEtherBalance("userKeys/" + user)
	if err != nil {
		userError(w, user, "", err)
		return
	}
	userSuccess(w, user, fmt.Sprintf("%d", bal))
}

// GetUserBalance - gets tokens and ether balances
func GetUserBalance(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if !validateUser(w, user) {
		userError(w, user, "", fmt.Errorf("invalid user"))
		return
	}
	etherBal, err := getUserEtherBalance("userKeys/" + user)
	if err != nil {
		userError(w, user, "", err)
		return
	}
	tokens, err := listTokens()
	if err != nil {
		adminError(w, "GetUserBalance", err)
		return
	}
	//fmt.Println(tokens)
	tokenBalances, err := getUserTokenBalances(user, tokens)

	blr := balanceListRec{Status: "OK", Result: tokenBalances}
	blr.Result = append(blr.Result, balance{Token: "ETH", Balance: etherBal})
	json.NewEncoder(w).Encode(blr)
}

func totalEtherBalance() (bal *big.Int, err error) {
	files, err := ioutil.ReadDir("userKeys/")
	if err != nil {
		return
	}
	bal = big.NewInt(0)
	etherBal := big.NewInt(0)
	for _, file := range files {
		etherBal, err = getUserEtherBalance("userKeys/" + file.Name())
		if err != nil {
			return
		}
		bal = big.NewInt(0).Add(bal, etherBal)
	}
	// Add the pool
	// etherBal, err = getUserEtherBalance("adminKeys/pool")
	// if err != nil {
	// 	return
	// }
	// bal = big.NewInt(0).Add(bal, etherBal)
	return
}

func getTotalUserBalances() (bals []balance, err error) {
	tokens, err := listTokens()
	if err != nil {
		return
	}
	tokenBal := big.NewInt(0)
	for _, tkn := range tokens {
		tokenBal, err = getTotalTokenBalance(tkn.Address)
		if err != nil {
			fmt.Println(2)
			return
		}
		bals = append(bals, balance{Token: tkn.Name, Balance: tokenBal})
	}
	tokenBal, err = totalEtherBalance()
	bals = append(bals, balance{Token: "ETH", Balance: tokenBal})
	return
}

// GetTotalUserBalances gets an array of total holdings per token
func GetTotalUserBalances(w http.ResponseWriter, r *http.Request) {
	bals, err := getTotalUserBalances()
	if err != nil {
		adminError(w, "GetTotalUserBalances", err)
		return
	}
	br := balanceListRec{Status: "OK", Result: bals}
	json.NewEncoder(w).Encode(br)
}
