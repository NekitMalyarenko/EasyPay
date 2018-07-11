package db

import (
	"types"
	"upper.io/db.v3"
	"errors"
)

type Verification string


const (
	verificationTable = "verification"
)


func (*Verification) CreateVerification(phoneNumber string, code int64,time string) error {
	_, err := currentInstance.instance.
		InsertInto(verificationTable).Values(types.Verification{
			PhoneNumber      : phoneNumber,
			IsVerified       : false,
			VerificationCode : code,
			StartTime        : time,
	}).Exec()
	return err
}


func (*Verification) GetVerification(phoneNumber string) (*types.Verification) {
	res := currentInstance.instance.Collection(verificationTable).Find(db.Cond{
		"phone_number" : phoneNumber,
	})

	var verification *types.Verification
	if res.Next(&verification) {
		return verification
	} else {
		return nil
	}
}


func (*Verification) Verify(phoneNumber string) error {
	res := currentInstance.instance.Collection(verificationTable).Find(db.Cond{
		"phone_number" : phoneNumber,
	})

	var verification types.Verification
	if res.Next(&verification) {
		return res.Update(map[string]interface{}{
			"is_verified" : true,
		})
	} else {
		return errors.New("no verification data")
	}
}


func (*Verification) IsVerified(phoneNumber string) (bool, error) {
	res := currentInstance.instance.Collection(verificationTable).Find(db.Cond{
		"phone_number" : phoneNumber,
	})

	var verification types.Verification
	if res.Next(&verification) {
		return verification.IsVerified, nil
	} else {
		return false, errors.New("no verification data")
	}
}


func (*Verification) DeleteVerification(phoneNumber string) error{
	return currentInstance.instance.Collection(verificationTable).Find(db.Cond{
		"phone_number" : phoneNumber}).Delete()
}
