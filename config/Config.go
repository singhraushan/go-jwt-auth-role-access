package config

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	username = "postgres"
	password = "admin"
	hostname = "127.0.0.1"
	port     = "5432"
	dbname   = "gopgTest"
	sslmode  = "disable"
	timeZone = "Asia/Kolkata"
)

//DB can use to connect DB from anywhere.
var DB *gorm.DB

func init() {
	fmt.Println("Config init")
	DB = getDBConnection()
}

func dsn() string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s TimeZone=%s", username, password, dbname, hostname, port, sslmode, timeZone)
}
func getDBConnection() *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}
