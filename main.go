package main

import (
	"github.com/joho/godotenv"
	"github.com/sitetester/ethereum-blockchain-parser/src/migration"
	"github.com/sitetester/ethereum-blockchain-parser/src/service/eth"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	user := getDotEnvVariable("MYSQL_USER_NAME")
	password := getDotEnvVariable("MYSQL_USER_PASSWORD")
	dsn := user + ":" + password + "@tcp(127.0.0.1:3306)/go_eth_blockchain?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect database!")
	}

	migration.RunMigrations(db)

	var importManager eth.ImportManager
	importManager.ManageImport(db)
}

func getDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
