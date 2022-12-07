package endpoints

import (
	"backend/blockchain"
	"fmt"
	"net/http"
)

func Wallets(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == "POST" {
		blockchain.CreateWallet()
	}
}
