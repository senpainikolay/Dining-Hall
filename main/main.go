package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

var Menu = order.GetFoods()
var conf = order.GetConf()
var orderId = order.OrderId{Id: 0}

func main() {

	rand.Seed(time.Now().UnixMilli())
	tables = order.GetTables(NumberOfTables)
	waiters = order.GetWaiters(NumberOfWaiters)
	rating := order.GetRating()

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")
	r.HandleFunc("/v2/order", PostClientOrders).Methods("POST")

	RegisterAtOM(conf.OMAddress, conf.RestaurantId, conf.LocalAddress, conf.RestaurantName, rating.Score)

	go func() {
		for {
			tables.OccupyTables()
			time.Sleep(order.TIME_UNIT * 80 * time.Millisecond)
		}
	}()
	for i := 0; i < NumberOfWaiters; i++ {
		idx := i
		go func() { waiters.Waiters[idx].Work(tables, &orderId, rating, conf.KitchenAddress, Menu) }()
	}

	http.ListenAndServe(":"+conf.Port, r)

}

func RegisterAtOM(address string, resId int, localAd string, resName string, rating float64) {
	postBody, _ := json.Marshal(order.RegisterPayload{
		RestaurantId: resId,
		Address:      localAd,
		Name:         resName,
		MenuItems:    len(Menu.Foods),
		Menu:         Menu.Foods,
		Rating:       rating,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+address+"/register", "application/json", responseBody)
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

func PostKitchenOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var ord order.Payload
	err := json.NewDecoder(r.Body).Decode(&ord)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	fmt.Fprint(w, "Succesfully recieved to Dining Hall")

	if ord.WaiterId == -1 || ord.TableId == -1 {
		// LOGIC TO BE IMPLEMENTERD ON CLIENT ORDERS STACKING
		log.Printf("CLIENT's ORDER ID %v at %v DONE \n ", ord.OrderId, conf.LocalAddress)
	} else {
		waiters.Waiters[ord.WaiterId-1].OrdersToServe <- ord
	}

}

func PostClientOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var clientOrd order.ClientOrder
	err := json.NewDecoder(r.Body).Decode(&clientOrd)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	orderId.Mutex.Lock()
	orderId.Id += 1
	temp := orderId.Id
	orderId.Mutex.Unlock()

	var ord = order.Order{
		OrderId:    temp,
		Items:      clientOrd.Items,
		Priority:   clientOrd.Priority,
		MaxWait:    clientOrd.MaxWait,
		PickUpTime: clientOrd.CreatedTime,
		TableId:    -1,
		WaiterId:   -1,
	}

	order.SendOrder(&ord, conf.KitchenAddress)

	clientOrdResponse := order.ClientOrderResponse{
		RestaurantId:         conf.RestaurantId,
		OrderId:              temp,
		EstimatedWaitingTime: 0,
		CreatedTime:          clientOrd.CreatedTime,
		RegisteredTime:       time.Now().UnixMilli(),
	}

	resp, _ := json.Marshal(&clientOrdResponse)
	fmt.Fprint(w, string(resp))

}
