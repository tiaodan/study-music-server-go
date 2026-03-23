package utils

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
