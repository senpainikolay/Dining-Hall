package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/senpainikolay/Dining-Hall/order"
)

func PostKitchenOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var ord order.Payload
	err := json.NewDecoder(r.Body).Decode(&ord)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	fmt.Fprint(w, "Succesfully recieved to Dining Hall")

	log.Printf("Order id %v succesfully recieved from Kitchen", ord.OrderId)

}

const (
	NumberOfTables  = 10
	NumberOfWaiters = 5
)

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")
	go http.ListenAndServe(":8080", r)

	var tables []order.Table

	for i := 1; i <= NumberOfTables; i++ {
		tables = append(tables,
			order.Table{
				TableId:         i,
				Free:            true,
				ReadyToOrder:    false,
				WaitingForOrder: false,
				OrderRecieving:  make(chan order.Order, 1),
			})
	}

	var waiters []order.Waiter

	for i := 1; i <= NumberOfWaiters; i++ {
		waiters = append(waiters,
			order.Waiter{
				WaiterId:        i,
				OrdersToRecieve: make(chan order.Order),
				OrdersToServe:   make(chan order.Order),
			})
	}

	orderId := order.OrderId{Id: 0}
	/*
		sendOrder := func(ord *order.Order) {
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
	*/

	go func() {
		for {
			order.OccupyTables(tables, NumberOfTables)
			time.Sleep(2 * time.Second)
		}
	}()

	for {

		go func() {

			for _, waiter := range waiters {

				select {
				case PostOrder := <-waiter.OrdersToRecieve:
					//sendOrder(&PostOrder)
					log.Printf(" Send to kitchen: %+v", PostOrder)

				case ServeOrder := <-waiter.OrdersToServe:
					log.Printf("%+v Serving", ServeOrder)

				default:
					time.Sleep(100 * time.Millisecond)
					//	order.OccupyTables(tables, NumberOfTables)
					waiter.PickUpOrder(tables, &orderId)

				}

			}
		}()

		time.Sleep(400 * time.Millisecond)

	}

}
