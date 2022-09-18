package order

import (
	"log"
	"time"
)

type Waiter struct {
	WaiterId        int           `json:"waiter_id"`
	PickUpTime      time.Duration `json:"pick_up_time"`
	OrdersToRecieve chan Order
	OrdersToServe   chan Order
}

func (w *Waiter) PickUpOrder(tables []Table, orderId *OrderId) {
	for i := 0; i < len(tables); i++ {
		idx := i
		go func() {
			tables[idx].Mutex.Lock()
			if tables[idx].ReadyToOrder == true {
				ord := GetRandomOrder(orderId)
				log.Printf("Waiter id %v picked up at table id %v", w.WaiterId, tables[idx].TableId)
				w.OrdersToRecieve <- *ord
				tables[idx].ReadyToOrder = false
				tables[idx].WaitingForOrder = true
			}
			tables[idx].Mutex.Unlock()
		}()

	}

}
