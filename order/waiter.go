package order

import (
	"log"
	"time"
)

type Waiter struct {
	WaiterId        int   `json:"waiter_id"`
	PickUpTime      int64 `json:"pick_up_time"`
	OrdersToRecieve chan Order
	OrdersToServe   chan Order
}

func (w *Waiter) PickUpOrder(tables []Table, orderId *OrderId) {
	for i := 0; i < len(tables); i++ {
		idx := i
		go func() {
			tables[idx].Mutex.Lock()
			defer tables[idx].Mutex.Unlock()
			if tables[idx].ReadyToOrder == true {
				ord := GetRandomOrder(orderId)
				ord.TableId = tables[idx].TableId
				ord.WaiterId = w.WaiterId
				ord.PickUpTime = time.Now().UnixNano()
				log.Printf("Waiter id %v picked up at table id %v", w.WaiterId, tables[idx].TableId)
				w.OrdersToRecieve <- *ord
				time.Sleep(200 * time.Millisecond)
				tables[idx].ReadyToOrder = false
				tables[idx].WaitingForOrder = true
			}

		}()

	}

}
