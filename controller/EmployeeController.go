package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/singhraushan/go-jwt-auth-role-access/entities"
	"github.com/singhraushan/go-jwt-auth-role-access/service"
)

//CreateCustomer validate input data and pass to service layer
func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--CreateCustomer start.")
	customer, err := parseValidateCustomer(w, r)
	if err != nil {
		return
	}
	log.Printf("conroller--CreateCustomer customer:%v", customer)
	w.Write([]byte(service.CreateCustomer(&customer)))
	log.Println("conroller--CreateCustomer end.")
}

//OpenAccount validate account input data and pass to service layer
func OpenAccount(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--OpenAccount start.")
	account, err := parseValidateAccount(w, r)
	if err != nil {
		return
	}
	w.Write([]byte(service.OpenAccount(&account)))
	log.Println("conroller--OpenAccount end.")
}
func UpdateKYC(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--UpdateKYC start.")
	customer, err := parseValidateCustomer(w, r)
	if err != nil {
		return
	}
	w.Write([]byte(service.UpdateKYC(&customer)))
	log.Println("conroller--UpdateKYC end.")
}
func GetCustomerDetails(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--GetCustomerDetails start.")
	s := strings.Split(r.URL.Path, "/")
	log.Println("split path value:", s[len(s)-1])
	custIDStr := s[len(s)-1]
	custID, err := strconv.Atoi(custIDStr)
	if err != nil {
		log.Println("controller---GetCustomerDetails invalid custId:", err)
		w.Write([]byte("Invalid Customer Id:" + custIDStr))
		return
	}
	w.Write([]byte(service.GetCustomerDetails(custID)))
	log.Println("conroller--GetCustomerDetails end.")
}
func DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--DeleteCustomer start.")
	s := strings.Split(r.URL.Path, "/")
	custIDStr := s[len(s)-1]
	custID, err := strconv.Atoi(custIDStr)
	if err != nil {
		log.Println("controller---DeleteCustomer invalid custId:", err)
		w.Write([]byte("Invalid Customer Id:" + custIDStr))
		return
	}
	w.Write([]byte(service.DeleteCustomer(custID)))
	log.Println("conroller--DeleteCustomer end.")
}
func GetBalance(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--GetBalance start.")
	accNo := r.Header.Get("accountNumber")

	accountNumber, err := strconv.Atoi(accNo)
	if err != nil {
		log.Println("controller---GetBalance invalid accountNumber:", err)
		w.Write([]byte("Invalid accountNumber:" + accNo))
		return
	}
	w.Write([]byte(service.GetBalance(accountNumber)))
	log.Println("conroller--GetBalance end.")
}
func MoneyTransfer(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--MoneyTransfer start.")
	fromAccNo := r.Header.Get("fromAccountNumber")
	fromAccountNumber, err := strconv.Atoi(fromAccNo)
	if err != nil {
		log.Println("controller---MoneyTransfer invalid from accountNumber:", err)
		w.Write([]byte("Invalid From accountNumber:" + fromAccNo))
		return
	}
	toAccNo := r.Header.Get("toAccountNumber")
	toAccountNumber, err := strconv.Atoi(toAccNo)
	if err != nil {
		log.Println("controller---MoneyTransfer invalid to accountNumber:", err)
		w.Write([]byte("Invalid To accountNumber:" + toAccNo))
		return
	}

	amountStr := r.Header.Get("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		log.Println("controller---MoneyTransfer invalid amount:", err)
		w.Write([]byte("Invalid To amount:" + amountStr))
		return
	}

	w.Write([]byte(service.MoneyTransfer(fromAccountNumber, toAccountNumber, amount)))
	log.Println("conroller--MoneyTransfer end.")
}
func PdfStatement(w http.ResponseWriter, r *http.Request) {
	log.Println("conroller--PdfStatement start.")
	accNo := r.Header.Get("accountNumber")

	accountNumber, err := strconv.Atoi(accNo)
	if err != nil {
		log.Println("controller---PdfStatement invalid accountNumber:", err)
		w.Write([]byte("Invalid accountNumber:" + accNo))
		return
	}
	fromDate := r.Header.Get("fromDate")
	toDate := r.Header.Get("toDate")

	w.Write([]byte(service.PdfStatement(accountNumber, fromDate, toDate)))
	log.Println("conroller--PdfStatement end.")
}
func parseValidateCustomer(w http.ResponseWriter, r *http.Request) (entities.Customer, error) {
	var customer entities.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		log.Println("controller:Error while parsing Customer", err)
		w.WriteHeader(http.StatusBadRequest)
		return customer, err
	}
	if validErrs := customer.Validate(); len(validErrs) > 0 {
		log.Println("controller:Validation error:", validErrs)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validErrs)
	}
	return customer, nil
}
func parseValidateAccount(w http.ResponseWriter, r *http.Request) (entities.Account, error) {
	var ac entities.Account
	err := json.NewDecoder(r.Body).Decode(&ac)
	if err != nil {
		log.Println("controller:Error while parsing Account", err)
		w.WriteHeader(http.StatusBadRequest)
		return ac, err
	}
	log.Println("controller:parseValidateAccount ", ac)
	if validErrs := ac.Validate(); len(validErrs) > 0 {
		log.Println("controller:Validation error:", validErrs)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validErrs)
	}
	return ac, nil
}
