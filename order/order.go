package order

import (
	"math/rand"
	"sync"
)

type Order struct {
	OrderId    int     `json:"order_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
	PickUpTime int64   `json:"pick_up_time"`
}

type OrderId struct {
	Id    int
	Mutex sync.Mutex
}

func GetRandomOrder(orderId *OrderId) *Order {
	menu := getFoods()
	orderId.Mutex.Lock()
	orderId.Id += 1
	temp := orderId.Id
	orderId.Mutex.Unlock()
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
