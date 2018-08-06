package local_storage

import (
	"types"
	"time"
)


type Payments struct {
	data []*types.PaymentData
}


const (
	defaultTTL = 30 * time.Minute
)


func (payments *Payments) PutPayment(payment *types.PaymentData) {
	payments.data = append(payments.data, payment)
}


func (payments *Payments) GetAllPayments()[]*types.PaymentData {
	return payments.data
}


func (payments *Payments) DeletePayment(index int) {
	payments.data = append(payments.data[:index], payments.data[index+1:]...)
}