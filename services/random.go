package services

import (
	"math/rand"
	"time"
)


func GetRandom(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	return rand.Int63n(max - min) + min
}