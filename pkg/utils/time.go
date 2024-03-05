package utils

import (
	"fmt"
	"time"
)

func GetOrderIdTime() (orderId string) {
	currentTime := time.Now().Nanosecond()
	orderId = fmt.Sprintf("%d", currentTime)
	return
}
