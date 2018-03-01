package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/DaveAppleton/etherUtils"
	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/DaveAppleton/etherdb"
	"github.com/ethereum/go-ethereum/common"
)

// var db *sql.DB

// func init() {
// 	var err error
// 	db, err = sql.Open("postgres", "host=localhost user=token password='erc20' port=5432	sslmode=disable dbname=tokenDB connect_timeout=10")
// 	if err != nil {
// 		fmt.Println("open token database : ", err)
// 		log.Fatal("open token database : ", err)
// 	}
// 	fmt.Println("database is open")
// 	//defer db.Close()
// }

type response struct {
	Status  string
	Error   string
	Message string
}

type tokenListRec struct {
	Status string
	Result []etherdb.Token
}

type balance struct {
	Token   string
	Balance *big.Int
}

type balanceListRec struct {
	Status string
	Result []balance
}

type userBalances struct {
	user     string
	balances []balance
}

func adminError(w http.ResponseWriter, msg string, err error) {
	resp := response{"ERROR", err.Error(), msg}
	json.NewEncoder(w).Encode(resp)
}

// AddTokenToSystem adds a token to the system given its address
func AddTokenToSystem(w http.ResponseWriter, r *http.Request) {
	var err error
	client, err = getClient()
	if err != nil {
		//clientError(w, err)
		return
	}

	addressStr := r.FormValue("address")

	tokens, err := getTokenData("", addressStr)
	if err != nil {
		adminError(w, "addTokenToSystem", err)
		return
	}
	if len(tokens) > 0 {
		adminError(w, "addTokenToSystem", errors.New(tokens[0].Name+" already added"))
		return
	}
	address := common.HexToAddress(addressStr)
	tokenContract, err := NewERC20(address, client)
	if err != nil {
		adminError(w, "addTokenToSystem (contract)", err)
		return
	}
	name, err := tokenContract.Name(nil)
	if err != nil {
		adminError(w, "addTokenToSystem (Name)", err)
		return
	}
	decimals, err := tokenContract.Decimals(nil)
	if err != nil {
		decimals = 0
	}
	symbol, err := tokenContract.Symbol(nil)
	if err != nil {
		adminError(w, "addTokenToSystem (Symbol)", err)
		return
	}
	fmt.Println(name, addressStr, symbol)
	tokenData := etherdb.Token{Name: name, Address: addressStr, Decimals: decimals, Symbol: symbol}
	err = tokenData.Add()
	if err != nil {
		adminError(w, "addTokenToSystem (DB)", err)
		return
	}
	err = json.NewEncoder(w).Encode(tokenData)
	if err != nil {
		adminError(w, "addTokenToSystem (result)", err)
	}
	json.NewEncoder(w).Encode(tokens)
}

func listTokens() (tokens []etherdb.Token, err error) {
	tokens, err = etherdb.GetAllTokens()
	return
}

// ListTokens returns an array of tekens known to the system
func ListTokens(w http.ResponseWriter, r *http.Request) {
	tokens, err := etherdb.GetAllTokens()
	if err != nil {
		adminError(w, "listTokens", err)
		return
	}
	var tokensRec tokenListRec
	tokensRec.Status = "OK"
	tokensRec.Result = tokens
	json.NewEncoder(w).Encode(tokens)
}

func getTokenData(symbol string, address string) (tokens []etherdb.Token, err error) {
	searchToken := etherdb.Token{Symbol: symbol, Address: address}
	tokens, err = searchToken.FindAll()
	return
}

// func getTokenDataFromChain(address string) {
// 	client, err := getClient()
// 	if err != nil {
// 		return
// 	}
// 	tknObj, err := NewERC20(common.HexToAddress(address), client)
// 	if err != nil {
// 		return
// 	}

// }

// GetTokenInfo retrieves info about a token from the database matching address or symbol
func GetTokenInfo(w http.ResponseWriter, r *http.Request) {
	var tokens tokenListRec
	var err error
	symbol := r.FormValue("symbol")
	address := r.FormValue("address")
	tokens.Result, err = getTokenData(symbol, address)
	if err != nil {
		adminError(w, "getTokenInfo", err)
		return
	}
	json.NewEncoder(w).Encode(tokens)
}

func getUserTokenBalances(user string, tokens []etherdb.Token) (bals []balance, err error) {
	var tknObj *ERC20
	client, err := getClient()
	if err != nil {
		return
	}
	userAddress, err := userAddress(user)
	if err != nil {
		return
	}
	fmt.Println(tokens)
	for _, tkn := range tokens {

		tknObj, err = NewERC20(common.HexToAddress(tkn.Address), client)
		if err != nil {
			return
		}
		bal := balance{Token: tkn.Name}

		bal.Balance, err = tknObj.BalanceOf(nil, userAddress)
		if err != nil {
			return
		}
		if bal.Balance.Cmp(big.NewInt(0)) != 0 {
			bals = append(bals, bal)
		}
	}
	return
}

// GetTokenBalance gets the balance of a specific token for one user
func GetTokenBalance(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	address := r.FormValue("address")
	user := r.FormValue("user")
	tokens, err := getTokenData(symbol, address)
	if err != nil {
		adminError(w, "getTokenBalance", err)
		return
	}
	if len(tokens) != 1 {
		adminError(w, "getTokenBalance", fmt.Errorf("%d tokens found with symbol '%s' or address '%s'", len(tokens), symbol, address))
		return
	}
	balances, err := getUserTokenBalances(user, tokens)
	if err != nil {
		adminError(w, "getTokenBalance", err)
		return
	}
	blr := balanceListRec{Status: "OK"}
	if len(balances) != 1 {
		blr.Result = []balance{balance{tokens[0].Name, big.NewInt(0)}}
	} else {
		blr.Result = []balance{balance{balances[0].Token, balances[0].Balance}}
	}
	json.NewEncoder(w).Encode(blr)
}

// get total balance for one token
func getTotalTokenBalance(tokenAddress string) (bal *big.Int, err error) {
	client, err := getClient()
	if err != nil {
		return
	}
	address := common.HexToAddress(tokenAddress)
	tokenObj, err := NewERC20(address, client)
	if err != nil {
		return
	}
	files, err := ioutil.ReadDir("userKeys/")
	if err != nil {
		return
	}
	bal = big.NewInt(0)
	tokenBal := big.NewInt(0)
	var userAddr common.Address
	for _, file := range files {

		userAddr, err = userAddress(file.Name())
		if err != nil {
			return
		}
		tokenBal, err = tokenObj.BalanceOf(nil, userAddr)
		if err != nil {
			return
		}
		bal = big.NewInt(0).Add(bal, tokenBal)
	}
	// // Add the pool

	// userAddr, err = userAddress("adminKeys/pool")
	// if err != nil {
	// 	return
	// }
	//tokenBal, err = tokenObj.BalanceOf(nil, userAddr)
	if err != nil {
		return
	}
	bal = big.NewInt(0).Add(bal, tokenBal)
	return
}

// GetTotalTokenBalances gets the sum of a specific token for all users + pool
func GetTotalTokenBalances(w http.ResponseWriter, r *http.Request) {
	symbol := r.FormValue("symbol")
	address := r.FormValue("address")
	tokens, err := getTokenData(symbol, address)
	if err != nil {
		adminError(w, "getTotalTokenBalances", err)
		return
	}
	if len(tokens) != 1 {
		adminError(w, "getTotalTokenBalances", fmt.Errorf("%d tokens found with symbol '%s' or address '%s'", len(tokens), symbol, address))
		return
	}

	bal, err := getTotalTokenBalance(tokens[0].Address)
	if err != nil {
		adminError(w, "getTotalTokenBalances", err)
		return
	}
	blr := balanceListRec{Status: "OK", Result: []balance{balance{Token: tokens[0].Name, Balance: bal}}}
	json.NewEncoder(w).Encode(blr)
}

// SendTokenFromUser - //
func SendTokenFromUser(w http.ResponseWriter, r *http.Request) {
	client, err := getClient()
	if err != nil {
		adminError(w, "SendTokenFromUser", err)
		return
	}
	userStr := r.FormValue("user")
	userKey := ethKeys.NewKey("userKeys/" + userStr)
	if userKey.LoadKey() != nil {
		adminError(w, "SendTokenFromUser", errors.New("User does not exist"))
		return
	}
	tx := userTx(userKey)

	destination := r.FormValue("destination")
	if len(destination) < 4 {
		adminError(w, "SendTokenFromUser", errors.New("invalid destination : "+destination))
	}
	symbol := r.FormValue("symbol")
	address := r.FormValue("address")
	tokens, err := getTokenData(symbol, address)
	if len(tokens) != 1 {
		adminError(w, "SendTokenFromUser", fmt.Errorf("%d tokens found", len(tokens)))
		return
	}

	gasLimitStr := r.FormValue("gasLimit")
	gasPriceStr := r.FormValue("gasPrice")
	var destinationAddress common.Address

	if strings.Compare(destination[:2], "0x") == 0 {
		destinationAddress = common.HexToAddress(destination)
	} else {
		destinationAddress, err = userAddress("userKeys/" + destination)
		if err != nil {
			adminError(w, "SendTokenFromUser", err)
			return
		}
	}
	valueStr := r.FormValue("value")

	value, ok := etherUtils.StrToDecimals(valueStr, int64(tokens[0].Decimals)) // must use decimals

	if !ok {
		adminError(w, "SendTokenFromUser", fmt.Errorf("Bad Number : %s", valueStr))
		return
	}
	gasLimit, err := strconv.Atoi(gasLimitStr)
	if err != nil {
		tx.GasLimit = uint64(gasLimit)
	}
	gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
	if ok {
		tx.GasPrice = gasPrice
	}

	tokenObj, err := NewERC20(common.HexToAddress(tokens[0].Address), client)
	if err != nil {
		adminError(w, "SendTokenFromUser", err)
		return
	}

	txn, err := tokenObj.Transfer(tx, destinationAddress, value)
	if err != nil {
		adminError(w, "SendTokenFromUser", err)
		return
	}
	csA := csAction{Hash: txn.Hash().Hex()}

	fmt.Println(csA)
	err = json.NewEncoder(w).Encode(csA)
	if err != nil {
		adminError(w, "SendEtherFromUser", err)
		fmt.Println(err)
	}
	return

}
