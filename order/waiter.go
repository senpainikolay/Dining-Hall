package order

import "time"

type Waiter struct {
	WaiterId        int           `json:"waiter_id"`
	PickUpTime      time.Duration `json:"pick_up_time"`
	OrdersToRecieve chan Order
	OrdersToServe   chan Order
}
