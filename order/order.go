package order

import (
	"math/rand"
)

type Order struct {
	OrderId  int     `json:"order_id"`
	Items    []int   `json:"items"`
	Priority int     `json:"priority"`
	MaxWait  float64 `json:"max_wait"`
}

func GetRandomOrder(orderId *int) *Order {
	menu := getFoods()
	*orderId += 1
	temp := *orderId
	items := make([]int, rand.Intn(10)+1)
	maxWaitInt := 0

	for i := range items {
		items[i] = rand.Intn(13) + 1
		if menu.Foods[items[i]-1].PreparationTime > maxWaitInt {
			maxWaitInt = menu.Foods[items[i]-1].PreparationTime
		}
	}
	priority := rand.Intn(4) + 1

	return &Order{
		OrderId:  temp,
		Items:    items,
		Priority: priority,
		MaxWait:  float64(maxWaitInt) * 1.3,
	}
}
