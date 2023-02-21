package mux

import (
	"IranStocksCrawler/system/router"
	"fmt"
	"net/http"

	gorillamux "github.com/gorilla/mux"
)

type muxRouter struct {
	router *gorillamux.Router
}

func NewMuxRouter() router.IRouter {
	return &muxRouter{router: gorillamux.NewRouter()}
}

func (r *muxRouter) HttpGet(uri string, f router.HandlerFunc) {

	r.router.HandleFunc(uri, f).Methods("GET")
}

func (r *muxRouter) HttpPost(uri string, f router.HandlerFunc) {

	r.router.HandleFunc(uri, f).Methods("POST")
}

func (r *muxRouter) HttpServe(port string) {

	fmt.Printf("We are listtening to the world on port %v by gurilla mux", port)
	http.ListenAndServe(":"+port, r.router)
}

func (r *muxRouter) GetType() router.RouterType {
	return router.MuxRouter
}
