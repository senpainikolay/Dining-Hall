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
	NumberOfWaiters = 5
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

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")

	tables := order.GetTables(NumberOfTables)
	waiters := order.GetWaiters(NumberOfWaiters)
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
			tables.OccupyTables()
			time.Sleep(5 * time.Second)
		}
	}()
	for i := 0; i < NumberOfWaiters; i++ {
		idx := i
		go func() { waiters.Waiters[idx].Work(tables, &orderId) }()
	}

	http.ListenAndServe(":8080", r)

}
