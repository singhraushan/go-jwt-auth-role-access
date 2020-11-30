package controller

import (
	"encoding/json"
	"net/http"

	"github.com/singhraushan/go-jwt-auth-role-access/encrypt"
	"github.com/singhraushan/go-jwt-auth-role-access/service"
)

type credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

//Login credential maching and token creation
func Login(w http.ResponseWriter, r *http.Request) {
	var creds credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	encryptedPassword, err := encrypt.Encrypt(creds.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	service.Login(w, encryptedPassword, creds.Username)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	service.Logout(w)
	w.Write([]byte("Successfully logged out"))
	w.WriteHeader(http.StatusOK)
}
