package local_storage

import (
	"sync"
	"types"
)

type localDBService struct {
	Payments
}


var (
	currentInstance *localDBService
	mu sync.Mutex
)


func GetInstance() (*localDBService) {
	mu.Lock()
	defer mu.Unlock()

	if currentInstance == nil {
		currentInstance = &localDBService{}
		currentInstance.data = make([]*types.PaymentData, 0)
	}

	return currentInstance
}
