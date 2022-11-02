package order

import (
	"math/rand"
	"sync"
)

type ClientOrder struct {
	Items       []int   `json:"items"`
	Priority    int     `json:"priority"`
	MaxWait     float64 `json:"max_wait"`
	CreatedTime int64   `json:"created_time"`
}

type Order struct {
	OrderId    int     `json:"order_id"`
	Items      []int   `json:"items"`
	Priority   int     `json:"priority"`
	MaxWait    float64 `json:"max_wait"`
	PickUpTime int64   `json:"pick_up_time"`
	TableId    int     `json:"table_id"`
	WaiterId   int     `json:"waiter_id"`
}

type ClientOrderResponse struct {
	RestaurantId         int     `json:"restaurant_id"`
	OrderId              int     `json:"order_id"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	CreatedTime          int64   `json:"created_time"`
	RegisteredTime       int64   `json:"registered_time"`
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

type EWTCalculation struct {
	B int `json:"B"`
	D int `json:"D"`
	E int `json:"E"`
}

func GetRandomOrder(orderId *OrderId, menu *Foods) *Order {

	orderId.Mutex.Lock()
	orderId.Id += 1
	temp := orderId.Id
	orderId.Mutex.Unlock()
	items := make([]int, rand.Intn(9)+1)
	maxWaitInt := 0

	for i := range items {
		items[i] = rand.Intn(len(menu.Foods)) + 1
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

type RestaurantRatingPayload struct {
	OrderId              int     `json:"order_id"`
	Rating               int     `json:"rating"`
	EstimatedWaitingTime float64 `json:"estimated_waiting_time"`
	WaitingTime          int     `json:"waiting_time"`
}

type RestaurantRatingResponse struct {
	RestaurantId        int     `json:"restaurant_id"`
	RestaurantAvgRating float64 `json:"restaurant_avg_rating"`
	PreparedOrders      int     `json:"prepared_orders"`
}
