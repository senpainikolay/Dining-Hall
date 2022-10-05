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

const (
	NumberOfTables  = 10
	NumberOfWaiters = 4
)

var waiters *order.Waiters
var tables *order.Tables

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
	log.Printf("%+v\n", ord)

	waiters.Waiters[ord.WaiterId-1].OrdersToServe <- ord

}

func main() {
	rand.Seed(time.Now().UnixMilli())
	tables = order.GetTables(NumberOfTables)
	waiters = order.GetWaiters(NumberOfWaiters)
	orderId := order.OrderId{Id: 0}
	rating := order.GetRating()

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")

	go func() {
		for {
			tables.OccupyTables()
			time.Sleep(order.TIME_UNIT * 80 * time.Millisecond)
		}
	}()
	for i := 0; i < NumberOfWaiters; i++ {
		idx := i
		go func() { waiters.Waiters[idx].Work(tables, &orderId, rating) }()
	}

	http.ListenAndServe(":8080", r)

}
