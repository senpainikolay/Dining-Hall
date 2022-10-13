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

type Waiter struct {
	WaiterId        int   `json:"waiter_id"`
	PickUpTime      int64 `json:"pick_up_time"`
	OrdersToRecieve chan Order
	OrdersToServe   chan Payload
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
				OrdersToServe:   make(chan Payload),
			})
	}

	return &waiters

}

func (w *Waiter) PickUpOrder(ts *Tables, orderId *OrderId) {

	for i := w.WaiterId - 1; i < len(ts.Tables)+w.WaiterId-1; i++ {
		idx := i
		idx = int(math.Mod(float64(idx), float64(len(ts.Tables)-1)))
		ts.Tables[idx].Mutex.Lock()
		if ts.Tables[idx].ReadyToOrder {
			ord := GetRandomOrder(orderId)
			ord.TableId = ts.Tables[idx].TableId
			ord.WaiterId = w.WaiterId
			ord.PickUpTime = time.Now().UnixMilli()
			//log.Printf("Waiter id %v picked up at table id %v", w.WaiterId, ts.Tables[idx].TableId)
			go func() { w.OrdersToRecieve <- *ord }()
			ts.Tables[idx].ReadyToOrder = false
			ts.Tables[idx].WaitingForOrder = true
			ts.Tables[idx].Mutex.Unlock()
			break
		}
		ts.Tables[idx].Mutex.Unlock()

	}

}

func (w *Waiter) Work(ts *Tables, orderId *OrderId, r *Rating, address string) {

	for {

		select {
		case ServeOrder := <-w.OrdersToServe:
			ok := (time.Now().UnixMilli() - ServeOrder.PickUpTime) / int64(TIME_UNIT)
			go func() {
				r.Mutex.Lock()
				r.Calculate(ServeOrder.MaxWait, float64(ok))
				r.Mutex.Unlock()
			}()
			//log.Printf("MAXWAIT: %v   THE TIME ?: %v", ServeOrder.MaxWait, ok)
			go func() {
				ts.Tables[ServeOrder.TableId-1].Mutex.Lock()
				ts.Tables[ServeOrder.TableId-1].WaitingForOrder = false
				time.Sleep(TIME_UNIT * 10 * time.Millisecond)
				ts.Tables[ServeOrder.TableId-1].Free = true
				ts.Tables[ServeOrder.TableId-1].Mutex.Unlock()
			}()
			//log.Printf("Waiter id %v serving table id %v with order id %v containing items: %+v \n", ServeOrder.WaiterId, ServeOrder.TableId, ServeOrder.OrderId, ServeOrder.Items)

		case PostOrder := <-w.OrdersToRecieve:
			sendOrder(&PostOrder, address)
			//log.Printf("Order id %v sent to kitchen: ", PostOrder.OrderId)

		default:
			time.Sleep(2 * TIME_UNIT * time.Millisecond)
			w.PickUpOrder(ts, orderId)

		}
	}

}

func sendOrder(ord *Order, address string) {
	postBody, _ := json.Marshal(*ord)
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+address+"/order", "application/json", responseBody)
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
