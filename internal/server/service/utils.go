package service

import (
	"math/rand"
	"time"
)

func getRandom() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(4-1+1) + 1
}
