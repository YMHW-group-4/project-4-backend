package API

import (
	"backend/API/endpoints"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

const PORD = ":2202"

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/wallets", endpoints.Wallets)
	http.HandleFunc("/transactions", endpoints.Transactions)
}

func StartServer() {
	handleRequests()
	bootHttpServer()
}

func bootHttpServer() {
	fmt.Println("About to listen on 8443. Go to http://127.0.0.1" + PORD)
	err := http.ListenAndServe(PORD, nil)
	if err != nil {
		log.Warn().Msg("Error while booting up the REST-API")
	}
}
