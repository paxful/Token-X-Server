package admin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/DaveAppleton/etherdb"
	"github.com/ethereum/go-ethereum/common"
)

type ueStruct struct {
	Status string
	Error  string
	User   string
	Result string
}

func userError(w http.ResponseWriter, user string, addr string, err error) {
	s := ueStruct{"ERROR", err.Error(), user, addr}
	json.NewEncoder(w).Encode(s)
}

func userSuccess(w http.ResponseWriter, user string, addr string) {
	s := ueStruct{"OK", "", user, addr}
	json.NewEncoder(w).Encode(s)
}

// If there is a particular requirement for user name...
func validateUser(w http.ResponseWriter, user string) bool {
	if len(user) > 0 {
		return true
	}
	userError(w, user, "", errors.New("Invalid User"))
	return false
}

// CreateUserAddress for user account
func CreateUserAddress(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if !validateUser(w, user) {
		return
	}
	key := ethKeys.NewKey("userKeys/" + user)
	if key.LoadKey() == nil {
		userError(w, user, key.PublicKeyAsHexString(), errors.New("User already exists"))
		return
	}
	if err := key.RestoreOrCreate(); err != nil {
		userError(w, user, "", err)
		return
	}
	if err := key.SaveKey(); err != nil {
		userError(w, user, "", err)
		return
	}
	acc := etherdb.Account{User: user, Address: key.PublicKeyAsHexString()}
	if err := acc.Add(); err != nil {
		log.Println(user, key.PublicKeyAsHexString(), err)
	}
	userSuccess(w, user, key.PublicKeyAsHexString())
}

func userAddress(user string) (addr common.Address, err error) {
	key := ethKeys.NewKey("userKeys/" + user)
	if key.LoadKey() != nil {
		err = errors.New("User does not exist")
		return
	}
	addr = key.PublicKey()
	return
}

// GetUserAddress of user
func GetUserAddress(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if !validateUser(w, user) {
		return
	}
	addr, err := userAddress(user)
	if err != nil {
		userError(w, user, "", errors.New("User does not exist"))
		return
	}
	userSuccess(w, user, addr.Hex())
}

// GetUserPrivateKey for user
func GetUserPrivateKey(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	if !validateUser(w, user) {
		return
	}
	key := ethKeys.NewKey("userKeys/" + user)
	if key.LoadKey() != nil {
		userError(w, user, "", errors.New("User does not exist"))
		return
	}
	userSuccess(w, user, fmt.Sprintf("0x%x", key.Key.D))
}

/// 0xcd6cf75d96b0fd6ae8f89952f37c9964701419e9 , 0x4f02d80841a0f4f899283462855491a11793fc33ed813c3307fd45a23b2a3277
