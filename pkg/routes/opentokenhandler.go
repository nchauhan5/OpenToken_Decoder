package routes

import (
	"fmt"
	"net/http"

	"bitbucket.org/imlabs/go_auth_service/pkg/logger"
	"bitbucket.org/imlabs/go_auth_service/pkg/opentoken"
	"github.com/gorilla/mux"
)

var log = logger.Logger

func addOpenTokenHandler(r *mux.Router) {
	r.HandleFunc("/auth", openTokenHandler).Methods("POST")
}

// openTokenHandler is used to handle a post request containing an OpenToken in its request payload/body
func openTokenHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/auth" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Println(w, "ParseForm() err: %v", err)
		return
	}

	// fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.ParseForm)
	token := r.Form.Get("opentoken")
	fmt.Println(w, "OpenToken = %s : ", token)

	var password = opentoken.AgentConfig.Password
	var cipherID = opentoken.AgentConfig.Ciphersuite
	opentoken.MakeStringBase64Compatible(&token)
	// fmt.Println(token)
	var decodedMsg, err = opentoken.DecodeToken(token, cipherID, password)
	if err != nil {
		fmt.Println(err)
	} else {
		opentoken.ProcessPayload(decodedMsg)
	}
}
