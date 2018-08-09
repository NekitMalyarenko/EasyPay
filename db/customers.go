package db

import (
	"types"
	"log"
	"upper.io/db.v3"
)


type Customers string


const (
	customersTable = "users"
)


func (*Customers) GetCustomer(phoneNumber string) (*types.Customer, error) {
	var user *types.Customer

	res := currentInstance.instance.Collection(customersTable).Find(db.Cond{"phone_number": phoneNumber})
	has, err := res.Exists()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if has {
		err := res.One(&user)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return user, nil
	} else {
		return nil, nil
	}
}


func (*Customers) GetCustomerById(id int64) (*types.Customer, error) {
	var user *types.Customer

	res := currentInstance.instance.Collection(customersTable).Find(db.Cond{"id": id})
	has, err := res.Exists()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if has {
		err := res.One(&user)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		return user, nil
	} else {
		return nil, nil
	}
}


func (*Customers) AddCustomer(user *types.Customer) error {
	_, err := currentInstance.instance.
		InsertInto(customersTable).Values(user).Exec()
	return err
}


func (*Customers) AddCustomerImage(id int64, image string) error {
	q := currentInstance.instance.Update("users").
		Set("image = ?", image).Where("id = ?", id)

	_, err := q.Exec()
	if err != nil {
		return err
	}

	return nil
}


func (db *Customers) HasCustomer(phoneNumber string) (bool, error) {
	user, err := db.GetCustomer(phoneNumber)
	if err != nil {
		return false, err
	} else if user != nil {
		return true, nil
	} else {
		return false, nil
	}
}