package router

import (
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type IRouter interface {
	HttpGet(uri string, f HandlerFunc)
	HttpPost(uri string, f HandlerFunc)
	HttpServe(port string)
	GetType() RouterType
}

type RouterType string

const (
	MuxRouter  RouterType = "mux"
	GinRouter  RouterType = "gin"
	EchoRouter RouterType = "echo"
)
