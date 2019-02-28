package httpserver

import (
	"fmt"
	"net/http"

	"bitbucket.org/imlabs/go_auth_service/pkg/logger"
	"bitbucket.org/imlabs/go_auth_service/pkg/routes"
)

var log = logger.Logger

// This init function just initialises a golang http server to run at localhost port:8080
func init() {
	fmt.Println("Inside httpserverstartup init() *************")

	fmt.Printf("Starting server ...\n")
	if err := http.ListenAndServe(":8080", routes.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
