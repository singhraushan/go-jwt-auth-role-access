package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const layout = "02 Jan 06 15:04 MST"

// Access invoke for role=employee user and check autherisation
func Access(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		okay, role := tokenValidation(w, r)
		if !okay {
			return
		}
		if role != "employee" {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised\n"))
			return
		}
		handler(w, r)
	}
}

//AdminAccess invoke for role=admin user and check autherisation
func AdminAccess(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		okay, role := tokenValidation(w, r)
		if !okay {
			return
		}
		if role != "admin" {
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised\n"))
			return
		}
		handler(w, r)
	}
}

//tokenValidation validating token
func tokenValidation(w http.ResponseWriter, r *http.Request) (bool, interface{}) {
	log.Println("middleware--tokenValidation start.")
	stringToken := r.Header.Get("authorization")
	claims := jwt.MapClaims{}
	jwtToken, err := jwt.ParseWithClaims(stringToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("secretKey"), nil //pass the key to decode
	})
	log.Println("jwtToken:", jwtToken, "---err:", err)
	if err != nil || !jwtToken.Valid {
		w.WriteHeader(401)
		w.Write([]byte("Invalid token: Unauthenticated\n"))
		return false, nil
	}
	log.Println("claims:", claims)
	expiryTime, _ := claims["expiryTime"]
	formatedExpiryTime, err := time.Parse(layout, fmt.Sprintf("%v", expiryTime))
	formatedExpiryTime.Format(layout)

	if formatedExpiryTime.Before(time.Now()) || formatedExpiryTime.Equal(time.Now()) {
		w.WriteHeader(401)
		w.Write([]byte("Token expired at:" + fmt.Sprintf("%v", expiryTime)))
		return false, nil
	}

	v, _ := claims["role"]
	log.Println("middleware--tokenValidation end.")
	return true, v
}
