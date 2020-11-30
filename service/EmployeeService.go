package service

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

	"github.com/shopspring/decimal"
	"github.com/singhraushan/go-jwt-auth-role-access/entities"
	"github.com/singhraushan/go-jwt-auth-role-access/repo"
	"github.com/singhraushan/go-jwt-auth-role-access/util"
)

//CreateCustomer will check in DB, whether this cudtomerid exist, if No then add new entry.
func CreateCustomer(c *entities.Customer) string {
	log.Printf("service--CreateCustomer start.customer:%v", c)
	repo.Create(c)
	log.Println("service--CreateCustomer end.")
	return "Customer created successfully"
}

//OpenAccount will check in DB, whether this cudtomerid exist, if No then add new entry.
func OpenAccount(ac *entities.Account) string {
	log.Println("service--OpenAccount start.")
	if !repo.IsCustomerExist(ac.CustomerID) {
		log.Println("service--OpenAccount custId:", ac.CustomerID, "already exist.")
		return "Customer Id " + fmt.Sprint(ac.CustomerID) + " is not present."
	}
	repo.Create(ac)
	log.Println("service--OpenAccount end.")
	return "Account opened successfully"
}

//UpdateKYC update customer details if present.
func UpdateKYC(c *entities.Customer) string {
	log.Println("service--UpdateKYC start.")
	if !repo.IsCustomerExist(c.ID) {
		log.Println("service--UpdateKYC custId:", c.ID, "already exist.")
		return "Customer Id " + fmt.Sprint(c.ID) + " is not present."
	}
	repo.UpdateCustomer(c, c.ID)
	log.Println("service--UpdateKYC end.")
	return "KYC updated successfully"
}

//GetCustomerDetails fetching customer and his/her account details if present.
func GetCustomerDetails(custID int) string {
	log.Println("service--GetCustomerDetails start.")
	var cust entities.Customer
	repo.GetRow(&cust, custID)
	if cust.ID == 0 {
		return "No data found. Customer Id does not exist."
	}
	var accounts []entities.Account
	repo.GetRows(&accounts)
	log.Println("accounts:", accounts)
	m := map[string]interface{}{}
	m["customer"] = cust
	m["account/s"] = accounts
	result, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	log.Println("service--GetCustomerDetails end.")
	return string(result)
}

//DeleteCustomer deleting customer and his/her account details if present.
func DeleteCustomer(custID int) string {
	log.Println("service--DeleteCustomer start.")
	var cust entities.Customer
	repo.GetRow(&cust, custID)
	if cust.ID == 0 {
		return "No data found. Customer Id does not exist."
	}
	repo.DeleteCustomerAndAccount(cust.ID)
	log.Println("service--DeleteCustomer end.")
	return "Customer details successfully deleted."
}

//GetBalance getting account balance for given accountNumber if present.
func GetBalance(accountNumber int) string {
	log.Println("service--GetBalance start.")
	var acc entities.Account
	repo.GetRow(&acc, accountNumber)
	if acc.CustomerID == 0 {
		return "No data found. Account Number does not exist."
	}
	log.Println("service--GetBalance end.")
	return "Account balance is:" + acc.Balance.String()
}

//MoneyTransfer transfer given account from 1st account to 2nd account.
func MoneyTransfer(fromAccountNumber, toAccountNumber int, amount float64) string {
	log.Println("service--MoneyTransfer start. Amount:", amount)
	var fromAcc entities.Account
	repo.GetRow(&fromAcc, fromAccountNumber)
	if fromAcc.CustomerID == 0 {
		return "No data found. From AccountNumber does not exist."
	}

	if v, _ := fromAcc.Balance.Float64(); v < amount {
		return fmt.Sprintf("Not enough amount to transfer. From acccount balance is:%f", v)
	}

	var toAcc entities.Account
	repo.GetRow(&toAcc, toAccountNumber)
	if toAcc.CustomerID == 0 {
		return "No data found. To AccountNumber does not exist."
	}
	if fromAccountNumber == toAccountNumber {
		return "Can't transfer within same account."
	}
	if amount <= 0 {
		return "Can't transfer negative/zero amount."
	}
	fromBlnc, _ := fromAcc.Balance.Float64()
	toBlnc, _ := toAcc.Balance.Float64()
	log.Println("service--MoneyTransfer fromBlnc:", fromBlnc, "toBlnc:", toBlnc, "fromBlnc-amount:", toFixed(fromBlnc-amount, 2), "toBlnc+amount:", toFixed(toBlnc+amount, 2))
	if !transferMoney(fromAccountNumber, toAccountNumber, toFixed(fromBlnc-amount, 2), toFixed(toBlnc+amount, 2), amount) {
		return "Could not able to transfer try later."
	}
	log.Println("service--MoneyTransfer end.")
	return fmt.Sprintf("From Account balance is:%v and To Account balance is:%v", toFixed(fromBlnc-amount, 2), toFixed(toBlnc+amount, 2))
}

//PdfStatement transfer given account from 1st account to 2nd account.
func PdfStatement(fromAccountNumber, fromDate, toDate interface{}) string {
	log.Println("service--PdfStatement start.")
	var acc entities.Account
	repo.GetRow(&acc, fromAccountNumber)
	if acc.CustomerID == 0 {
		return "AccountNumber does not exist."
	}
	var transactions []entities.Transaction
	repo.Statement(&transactions, fromAccountNumber, fromDate, toDate)
	if len(transactions) == 0 {
		return "No Transaction found."
	}
	log.Println("service--PdfStatement transactions:", transactions)
	util.GeneratePdfReport(&transactions)
	log.Println("service--PdfStatement end.")
	return "Report generated successfully."
}

func transferMoney(fromAccountNumber, toAccountNumber int, fromBalance, toBalance, amount float64) bool {

	if repo.UpdateAccount(fromAccountNumber, fromBalance) {
		if repo.UpdateAccount(toAccountNumber, toBalance) {
			//make Transaction
			t1 := entities.Transaction{
				FromAccountNumber: fromAccountNumber,
				ToAccountNumber:   toAccountNumber,
				AvailableAmount:   decimal.NewFromFloat(fromBalance),
				TranferAmount:     decimal.NewFromFloat(amount),
				Type:              "DEBIT",
			}
			repo.Create(&t1)
			t2 := entities.Transaction{
				FromAccountNumber: fromAccountNumber,
				ToAccountNumber:   toAccountNumber,
				AvailableAmount:   decimal.NewFromFloat(toBalance),
				TranferAmount:     decimal.NewFromFloat(amount),
				Type:              "CREDIT",
			}
			repo.Create(&t2)
			return true
		}
		repo.UpdateAccount(fromAccountNumber, toFixed(fromBalance+amount, 2))
		return false
	}
	return false
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
