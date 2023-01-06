package api

import (
	"context"
	"crypto/tls"
	"net/http"
)

type API struct {
	ctx     context.Context
	handler http.Handler
	config  *tls.Config
}

func NewAPI(port int) *API {
	return &API{
		ctx: context.Background(),
	}
}

//func (a *Api) handleRequests() {
//	http.HandleFunc()
//}

func (a *API) Start() error {
	return nil
}
