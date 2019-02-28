package routes

import (
	"fmt"
	"encoding/json"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"bitbucket.org/imlabs/go_auth_service/pkg/accounts"
)

func addUserAccountHandler(r *mux.Router) {
	r.HandleFunc("/user", userAccountHandler).Methods("POST")
}

func userAccountHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user accounts.UserModel
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	user.ID = bson.NewObjectId()
	if err := accounts.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, user)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
