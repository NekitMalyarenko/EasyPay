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


func (*Users) GetUser(email string) (*types.User, error) {
	var user *types.User

	res := currentInstance.instance.Collection(usersTable).Find(db.Cond{"email": email})
	err := res.One(&user)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return user, nil
}


func (*Users) AddUser(user *types.User) error {
	_, err := currentInstance.instance.
		InsertInto(usersTable).Values(user).Exec()
	return err
}


func (db *Users) HasUser(email string) (bool, error) {
	user, err := db.GetUser(email)
	if err != nil {
		return false, err
	} else if user != nil {
		return true, nil
	} else {
		return false, nil
	}
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



