package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

// Connect initializes the database connection
func Connect(pDbUser, pDbPass, pDbHost, pDbPort, pDbName string) error {
	log.Println("DataBase Connection (+)")
	defer log.Println("DataBase Connection (-)")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		pDbUser, pDbPass, pDbHost, pDbPort, pDbName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf(" Failed to connect to database: %v", err)
		return err
	}

	log.Println(" Database connected successfully")
	return nil
}
