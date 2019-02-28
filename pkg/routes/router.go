package routes

import (
	"github.com/gorilla/mux"
)

//NewRouter is main routing entry point
func NewRouter() *mux.Router {
	r := mux.NewRouter()

	addUserAccountHandler(r) // UserAccountHandler
	addOpenTokenHandler(r)   // OpenTokenHandler

	return r
}
