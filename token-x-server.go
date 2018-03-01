package main

import (
	"fmt"
	"net/http"

	"github.com/DaveAppleton/etherdb"

	"github.com/DaveAppleton/ether_go/ethKeys"

	"github.com/spf13/viper"

	"./admin"
)

func main() {
	paxful := ethKeys.NewKey("adminKeys/paxful")
	paxful.RestoreOrCreate()
	fmt.Println("Paxful account is at ", paxful.PublicKeyAsHexString())
	fmt.Println("Token-X-Server")
	viper.SetConfigName("config")                // name of config file (without extension)
	viper.AddConfigPath("/etc/Token-X-Server/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.Token-X-Server") // call multiple times to add many search paths
	viper.AddConfigPath(".")                     // optionally look for config in the working directory
	err := viper.ReadInConfig()                  // Find and read the config file
	if err != nil {                              // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
	etherdb.InitDB(viper.GetString("postgres_connect"))
	http.HandleFunc("/admin/createUserAddress", admin.CreateUserAddress)       // re-tested
	http.HandleFunc("/admin/getUserAddress", admin.GetUserAddress)             // re-tested
	http.HandleFunc("/admin/getUserPrivateKey", admin.GetUserPrivateKey)       // re-tested
	http.HandleFunc("/admin/getUserEtherBalance", admin.CheckUserEtherBalance) // tested
	http.HandleFunc("/admin/getUserBalance", admin.GetUserBalance)             // tested - eth + 2 x token
	http.HandleFunc("/admin/syncProgress", admin.SyncProgress)
	http.HandleFunc("/admin/getTokenInfo", admin.GetTokenInfo)
	http.HandleFunc("/admin/sendEtherFromUser", admin.SendEtherFromUser)         // tested
	http.HandleFunc("/admin/addTokenToSystem", admin.AddTokenToSystem)           // tested
	http.HandleFunc("/admin/listTokens", admin.ListTokens)                       // tested
	http.HandleFunc("/admin/getTotalUserBalances", admin.GetTotalUserBalances)   // tested
	http.HandleFunc("/admin/getTotalTokenBalances", admin.GetTotalTokenBalances) // tested
	http.HandleFunc("/admin/getTokenBalance", admin.GetTokenBalance)             // tested
	http.HandleFunc("/admin/getTxStatus", admin.GetTransactionStatus)
	http.HandleFunc("/admin/sendTokenFromUser", admin.SendTokenFromUser)
	port := viper.GetString("port")
	http.ListenAndServe(port, nil)
}
