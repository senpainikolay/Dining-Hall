package order

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

const TIME_UNIT = 100

func GetTimeUnit() int {
	return TIME_UNIT
}

type Waiter struct {
	WaiterId        int   `json:"waiter_id"`
	PickUpTime      int64 `json:"pick_up_time"`
	OrdersToRecieve chan Order
	OrdersToServe   chan Order
}

type Waiters struct {
	Waiters []Waiter
}

func GetWaiters(numberOfWaiters int) *Waiters {
	var waiters Waiters
	for i := 1; i <= numberOfWaiters; i++ {
		waiters.Waiters = append(waiters.Waiters,
			Waiter{
				WaiterId:        i,
				OrdersToRecieve: make(chan Order),
				OrdersToServe:   make(chan Order),
			})
	}

	return &waiters

}

func (w *Waiter) PickUpOrder(ts *Tables, orderId *OrderId) {
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()
	for i := w.WaiterId - 1; i < len(ts.Tables)+w.WaiterId-1; i++ {
		idx := i
		idx = int(math.Mod(float64(idx), float64(len(ts.Tables)-1)))
		if ts.Tables[idx].ReadyToOrder {
			ord := GetRandomOrder(orderId)
			ord.TableId = ts.Tables[idx].TableId
			ord.WaiterId = w.WaiterId
			ord.PickUpTime = time.Now().UnixNano()
			log.Printf("Waiter id %v picked up at table id %v", w.WaiterId, ts.Tables[idx].TableId)
			go func() { w.OrdersToRecieve <- *ord }()
			ts.Tables[idx].ReadyToOrder = false
			ts.Tables[idx].WaitingForOrder = true
			break
		}

	}

}

func (w *Waiter) Work(ts *Tables, orderId *OrderId) {

	for {

		select {
		case PostOrder := <-w.OrdersToRecieve:
			sendOrder(&PostOrder)
			log.Printf(" Send to kitchen: %+v", PostOrder)

		case ServeOrder := <-w.OrdersToServe:
			log.Printf("%+v Serving", ServeOrder)

		default:
			time.Sleep(4 * TIME_UNIT * time.Millisecond)
			w.PickUpOrder(ts, orderId)

		}
	}

}

func sendOrder(ord *Order) {
	postBody, _ := json.Marshal(*ord)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://localhost:8081/order", "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)
	log.Printf(sb)

}
