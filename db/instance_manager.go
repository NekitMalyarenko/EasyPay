package db

import (
	"upper.io/db.v3/postgresql"
	"log"
	"upper.io/db.v3/lib/sqlbuilder"
	"sync"
)


type dbService struct {
	instance sqlbuilder.Database
	Customers
	Sellers
	Verification
	Transactions
	Shops
}


var (
	currentInstance *dbService

	settings = postgresql.ConnectionURL{
		Host:     "ec2-54-247-98-162.eu-west-1.compute.amazonaws.com",
		Database: "df7qtm51mteljj",
		User:     "uzoysoozzaqazl",
		Password: "aa2f03f8e8b40ef7d893f95a91cd2d22a8d1e690ee6e6afbb254e5b0e4d43473",
		Options:map[string]string{"sslmode" : "require"},
	}

	mu sync.Mutex
)


func GetInstance() *dbService {
	mu.Lock()
	defer mu.Unlock()

	if currentInstance == nil {
		log.Println("-----New DB instance-----")
		newInstance, err := postgresql.Open(settings)
		if err != nil {
			log.Fatal(err)
		}

		currentInstance = &dbService{
			instance: newInstance,
		}
	}

	return currentInstance
}


func (db *dbService) CloseConnection() {
	db.instance.Close()
}