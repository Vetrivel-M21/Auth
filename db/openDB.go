package db

import (
	"auth/common"
	"auth/config"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB is the global database connection
var DB *gorm.DB

func OpenConnection() error {
	var cfg config.Config

	if err := common.LoadTOMLConfig("../dbconfig.toml", &cfg); err != nil {
		log.Println("Error while load toml (DOC001)", err)
		log.Fatal(err)
	}

	config.SecretKey = cfg.SECRET_KEY

	lErr := Connect(cfg.DB_USER, cfg.DB_PASS, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	if lErr != nil {
		log.Println("Connetion Failed (DOC002)", lErr)
		return lErr
	}
	return nil
}

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
