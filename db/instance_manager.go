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
	Ratings
}


var (
	currentInstance *dbService

	settings = postgresql.ConnectionURL{
		Host:     "ec2-79-125-12-48.eu-west-1.compute.amazonaws.com",
		Database: "dcpeust2o38fga",
		User:     "ixuneudikjaefb",
		Password: "e1119f3e0c0bdd0c5cad33dc8fb53470c1c69c0317743f03809725ec09b25f4a",
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