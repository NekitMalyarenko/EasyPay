package db

import (
	"types"
	"log"
	"upper.io/db.v3"
)

type Users string


const (
	usersTable = "users"
)


func (*Users) GetUser(phoneNumber string) *types.User {

	return nil
}


func (*Users) AddUser(user *types.User) error {
	_, err := currentInstance.instance.
		InsertInto(usersTable).Values(user).Exec()
	return err
}


func (*Users) HasUser(phoneNumber string) bool {
	var users []*types.User

	res := currentInstance.instance.Collection(usersTable).Find(db.Cond{"phone_number": phoneNumber})
	err := res.All(&users)
	if err != nil {
		log.Println(err)
		return false
	}

	return len(users) != 0
}


func (*Users) UserLogin(email, password string) bool {
	var users []*types.User

	res := currentInstance.instance.Collection(usersTable).Find(db.Cond{
		"email": email,
		"password" : password})

	err := res.All(&users)
	if err != nil {
		log.Println(err)
		return false
	}

	return len(users) != 0
}



