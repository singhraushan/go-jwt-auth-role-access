package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/singhraushan/go-jwt-auth-role-access/controller"
	"github.com/singhraushan/go-jwt-auth-role-access/middleware"
)

func main() {
	router := mux.NewRouter().StrictSlash(true) //if end with slice then only get result
	handleRequests(router)
	log.Fatal(http.ListenAndServe(":12345", router))
}

func handleRequests(router *mux.Router) {
	//any valid user
	router.HandleFunc("/signIn", controller.Login)
	router.HandleFunc("/signOut", controller.Logout)
	//only admin
	router.Handle("/addEmployee", middleware.AdminAccess(controller.AddNewEmployee)).Methods("POST")
	router.HandleFunc("/deleteEmployee", middleware.AdminAccess(controller.DeleteEmployee)).Methods("DELETE")
	//any employee
	router.HandleFunc("/createCustomer", middleware.Access(controller.CreateCustomer)).Methods("POST")
	router.HandleFunc("/customer/createAccount", middleware.Access(controller.OpenAccount)).Methods("POST") //saving,salary,salary,current
	//Link custome with account?Don't understand this requiremnt since openAccount link already available.
	router.HandleFunc("/customer/updateKYC", middleware.Access(controller.UpdateKYC)).Methods("PUT")
	router.HandleFunc("/customer/{id}", middleware.Access(controller.GetCustomerDetails)).Methods("GET")
	router.HandleFunc("/customer/{id}", middleware.Access(controller.DeleteCustomer)).Methods("DELETE")
	router.HandleFunc("/account/balance", middleware.Access(controller.GetBalance)).Methods("GET") //get account ID from header
	router.HandleFunc("/account/moneytranfer", middleware.Access(controller.MoneyTransfer)).Methods("PUT")
	router.HandleFunc("/account/statement", middleware.Access(controller.PdfStatement)).Methods("GET")
}
