package admin

import (
	"fmt"
	"math/big"

	"context"

	"github.com/DaveAppleton/ether_go/ethKeys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func sendEthereum(sender *ethKeys.AccountKey, recipient common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) (*types.Transaction, error) {
	//var ret interface{}
	//var zero interface{}

	ec, _ := getClient()

	nonce, err := ec.PendingNonceAt(context.TODO(), sender.PublicKey())
	if gasPrice == nil {
		gasPrice, err = ec.SuggestGasPrice(context.TODO())
		if err != nil {
			return nil, err
		}
	}
	fmt.Println("Nonce : ", nonce)
	fmt.Println("GasPrice : ", gasPrice)
	s := types.HomesteadSigner{}

	t := types.NewTransaction(nonce, recipient, amount, gasLimit, gasPrice, data)
	nt, err := types.SignTx(t, s, sender.GetKey())
	if err != nil {
		return nil, err
	}
	err = ec.SendTransaction(context.TODO(), nt)
	return nt, err
}
