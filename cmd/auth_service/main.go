package main

import (
	"fmt"

	_ "bitbucket.org/imlabs/go_auth_service/pkg/authttpserver"
	"bitbucket.org/imlabs/go_auth_service/pkg/logger"
)

var log = logger.Logger

func init() {
	fmt.Println("Inside main init() *************")
}

func main() {
	// var token = "T1RLAQLVVgI6nfAXif1wYQz-4Hoqqjpk-RCRhrYo_A3vfozy8DwQgX_iAAAgXtSyTiGFVbQGmJ7-USFFjaZYuPueXSr8Gl2W5APuFWw*"

	// // fmt.Println(token)
	// // var password = os.Getenv("password")
	// var password = opentoken.AgentConfig.Password
	// var cipherID = 2
	// opentoken.MakeStringBase64Compatible(&token)
	// // fmt.Println(token)
	// var decodedMsg, err = opentoken.DecodeToken(token, cipherID, password)
	// if err != nil {
	// 	os.Exit(1)
	// } else {
	// 	opentoken.ProcessPayload(decodedMsg)
	// }
}
