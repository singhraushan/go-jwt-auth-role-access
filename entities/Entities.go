package entities

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var availableAccountType = map[string]int{
	"SAVINGS": 1,
	"SALARY":  2,
	"LOAN":    3,
	"CURRENT": 4,
}

//UserDetails for user information
type UserDetails struct {
	Username string `gorm:"primarykey"`
	Password string
	Role     string
	Email    string
}

//Validate user's mandatory information
func (ud *UserDetails) Validate() url.Values {
	errs := url.Values{}

	if ud.Username == "" {
		errs.Add("Username", "The Username field is required!")
	}
	if ud.Password == "" {
		errs.Add("Password", "The Username field is required!")
	}
	if ud.Role == "" {
		errs.Add("Role", "The Role field is required!")
	}
	if ud.Email == "" {
		errs.Add("Email", "The Email field is required!")
	}
	return errs
}

//Customer information
type Customer struct {
	ID            int `gorm:"primary_key, AUTO_INCREMENT"`
	Name          string
	ContactNumber string
	Email         string
	Address       string
	DOB           string
}

//Validate customer's mandatory information
func (c *Customer) Validate() url.Values {
	errs := url.Values{}

	if c.Name == "" {
		errs.Add("Name", "The Name field is required!")
	}
	if c.ContactNumber == "" {
		errs.Add("ContactNumber", "The ContactNumber field is required!")
	}
	if c.Address == "" {
		errs.Add("Address", "Address information is required!")
	}
	return errs
}

//Account information
type Account struct {
	AccountNumber int `gorm:"PRIMARY_KEY;AUTO_INCREMENT;autoIncrement:true"`
	CustomerID    int
	Type          string
	Balance       decimal.Decimal `gorm:"type:numeric"`
}

//Validate customer's mandatory information
func (ac *Account) Validate() url.Values {
	errs := url.Values{}
	if ac.Type == "" {
		errs.Add("Type", "The Name field is required!")
	}
	if _, okay := availableAccountType[ac.Type]; !okay {
		accTypes, err := json.Marshal(availableAccountType)
		if err != nil {
			log.Println("entities---availableAccountType marsing error ", err)
			errs.Add("Type", "Invalid account type.")
		} else {
			errs.Add("Type", "The Account type should belongs to:"+string(accTypes))
		}
	}
	return errs
}

//Transaction information
type Transaction struct {
	gorm.Model
	FromAccountNumber int
	ToAccountNumber   int
	AvailableAmount   decimal.Decimal `gorm:"type:numeric"`
	TranferAmount     decimal.Decimal `gorm:"type:numeric"`
	Type              string
}
