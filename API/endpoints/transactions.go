package endpoints

import "net/http"

func Transactions(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		createTransaction()
	}

}

func createTransaction() {
	// TODO: call the code dat makes a transaction
}
