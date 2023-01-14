package endpoints

import (
	"fmt"
	"net/http"
)

func Transaction(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello from the server!\n")
}
