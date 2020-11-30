package service

import (
	"log"

	"github.com/singhraushan/go-jwt-auth-role-access/entities"
	"github.com/singhraushan/go-jwt-auth-role-access/repo"
)

//AddNewEmployee will check in DB, whether this username entry exist, if No then add new entry.
func AddNewEmployee(ud entities.UserDetails) string {
	log.Println("service--AddNewEmployee start.")
	if repo.IsUserExisting(ud.Username) {
		log.Println("service--AddNewEmployee username:", ud.Username, "already exist.")
		return "username: " + ud.Username + " already exist"
	}
	repo.Create(ud)
	log.Println("service--AddNewEmployee end.")
	return "Employee added successfully"
}

//DeleteEmployee will check in DB is this username entry exist, if Yes then delete.
func DeleteEmployee(ud entities.UserDetails) string {
	log.Println("service--DeleteEmployee start.")
	if !repo.IsUserExisting(ud.Username) {
		log.Println("service--DeleteEmployee: Employee does not exist")
		return "Employee does not exist"
	}
	repo.Delete(ud)
	log.Println("service--DeleteEmployee end.")
	return "Employee deleted successfully"
}
