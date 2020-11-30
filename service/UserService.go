package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/singhraushan/go-jwt-auth-role-access/entities"
	"github.com/singhraushan/go-jwt-auth-role-access/repo"
	"github.com/singhraushan/go-jwt-auth-role-access/util"
)

//Login credential maching and token creation
func Login(w http.ResponseWriter, encryptedPassword, username string) {
	log.Println("service--Login: encryptedPassword:", encryptedPassword, "username:", username)
	var ud entities.UserDetails
	repo.GetUserRow(&ud, username)

	log.Println("### DB data username:", ud.Username, "---password:", ud.Password, "len(ud.Username):", len(ud.Username))
	if len(ud.Username) == 0 || ud.Password != encryptedPassword {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid username/password"))
		return
	}
	//create token
	token, err := util.CreateToken(ud.Username, ud.Role)
	if err != nil {
		log.Println(err)
		return
	}
	json.NewEncoder(w).Encode(token)
}

//Logout session out
func Logout(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(0) //token pass as 0
}
