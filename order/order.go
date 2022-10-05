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
	PickUpTime int64   `json:"pick_up_time"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
}

type Payload struct {
	OrderId        int              `json:"order_id"`
	Items          []int            `json:"items"`
	Priority       int              `json:"priority"`
	MaxWait        float64          `json:"max_wait"`
	PickUpTime     int64            `json:"pick_up_time"`
	TableId        int              `json:"table_id"`
	WaiterId       int              `json:"waiter_id"`
	CookingTime    int64            `json:"cooking_time"`
	CookingDetails []CookingDetails `json:"cooking_details"`
}
type CookingDetails struct {
	CookId int `json:"cook_id"`
	FoodId int `json:"food_id"`
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
	items := make([]int, rand.Intn(5)+10)
	maxWaitInt := 0

	for i := range items {
		items[i] = rand.Intn(13) + 1
		if menu.Foods[items[i]-1].PreparationTime > maxWaitInt {
			maxWaitInt = menu.Foods[items[i]-1].PreparationTime
		}
	}
	priority := rand.Intn(5) + 1

	return &Order{
		OrderId:  temp,
		Items:    items,
		Priority: priority,
		MaxWait:  float64(maxWaitInt) * 1.3,
	}
}
