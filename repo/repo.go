package repo

import (
	"database/sql"
	"log"

	"github.com/singhraushan/go-jwt-auth-role-access/config"
	"github.com/singhraushan/go-jwt-auth-role-access/entities"
)

//Common for all table
func init() {
	//droping tables
	config.DB.Migrator().DropTable(&entities.UserDetails{})
	config.DB.Migrator().DropTable(&entities.Customer{})
	config.DB.Migrator().DropTable(&entities.Account{})
	config.DB.Migrator().DropTable(&entities.Transaction{})
	//creating tables
	config.DB.AutoMigrate(&entities.UserDetails{})
	config.DB.AutoMigrate(&entities.Customer{})
	config.DB.AutoMigrate(&entities.Account{})
	config.DB.AutoMigrate(&entities.Transaction{})

	ud := []entities.UserDetails{
		{Username: "raushan", Password: "rDTedc4=", Role: "admin", Email: "raushan@gmail.com"},    //password:admin
		{Username: "anish", Password: "qD3DcM/EQPY=", Role: "employee", Email: "anish@gmail.com"}, //password:employee
	}
	config.DB.Create(&ud)
}

//GetRow will search one entry of given struct based on primary key
func GetRow(t interface{}, pk interface{}) {
	config.DB.First(t, pk)
}

//GetRows will fetch all data of given struct table
func GetRows(t interface{}) {
	config.DB.Find(t)
}

//Create will insert new entry from DB
func Create(t interface{}) {
	config.DB.Create(t)
}

//Delete entry from DB
func Delete(t interface{}) {
	config.DB.Delete(t)
}

//Table Specific

//GetUserRow will search one entry of UserDetails based on primary key
func GetUserRow(ud *entities.UserDetails, pk interface{}) {
	config.DB.Model(&entities.UserDetails{}).Where("username", pk).Find(ud)
}

//IsUserExisting hit db and check UserDetails has entry for given username.
func IsUserExisting(username string) bool {
	var ud entities.UserDetails
	rows, err := config.DB.Model(&ud).Where("username", username).Rows()
	if err != nil || rows.Next() {
		return true
	}
	return false
}

//IsCustomerExist hit db and check Customer table has entry for given customerId.
func IsCustomerExist(custID int) bool {
	var c entities.Customer
	rows, err := config.DB.Model(&c).Where("id", custID).Rows()
	if err != nil || rows.Next() {
		return true
	}
	return false
}

//UpdateCustomer update customer with given data
func UpdateCustomer(c interface{}, id interface{}) {
	var typ entities.Customer
	config.DB.Model(&typ).Where("id", id).Updates(c)
}

//DeleteCustomerAndAccount delete cutomer and acount data for given cutomer id
func DeleteCustomerAndAccount(id interface{}) {
	var c entities.Customer
	config.DB.Where("id", id).Delete(&c)
	var acc entities.Account
	config.DB.Where("CustomerID", id).Delete(&acc)
}

//UpdateAccount update customer with given data
func UpdateAccount(acc interface{}, balance interface{}) bool {
	var typ entities.Account
	log.Println("repo--UpdateAccount balance: ", balance)
	config.DB.Model(&typ).Where("account_number", acc).Update("balance", balance)
	return true
}

//Statement ensure to fetch transactions
func Statement(type1 *[]entities.Transaction, acc, fromDate, toDate interface{}) {

	config.DB.Raw("SELECT * FROM transactions WHERE created_at >= @fromDate AND created_at <= @toDate AND (from_account_number = @acc AND type = @type1) OR (to_account_number = @acc AND type = @type2)",
		sql.Named("acc", acc), sql.Named("fromDate", fromDate), sql.Named("toDate", toDate), sql.Named("type1", "DEBIT"), sql.Named("type2", "CREDIT")).Find(type1)
}
