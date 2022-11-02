package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/senpainikolay/Dining-Hall/order"
)

const (
	NumberOfTables  = 6
	NumberOfWaiters = 3
)

var waiters *order.Waiters
var tables *order.Tables

var rating = order.GetRating()
var Menu = order.GetFoods()
var conf = order.GetConf()
var orderId = order.OrderId{Id: 0}

var PickUpController = order.OrderPickUpController{Mutex: sync.Mutex{}, SignalVar: 0}

var clientMap = order.GetClientMap()

func main() {

	rand.Seed(time.Now().UnixMilli())
	tables = order.GetTables(NumberOfTables)
	waiters = order.GetWaiters(NumberOfWaiters)

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")
	r.HandleFunc("/v2/order", PostClientOrders).Methods("POST")
	r.HandleFunc("/getOrderStatus", GetOrdersStatus).Methods("GET")
	r.HandleFunc("/v2/order/{id}", GetClientOrderDetails).Methods("GET")
	r.HandleFunc("/v2/rating", ClientRatingPost).Methods("POST")

	RegisterAtOM(rating.Score)

	go func() {
		for {
			tables.OccupyTables(conf.KitchenAddress, &PickUpController)
			time.Sleep(order.TIME_UNIT * time.Duration((rand.Intn(20) + 60)) * time.Millisecond)

		}
	}()
	for i := 0; i < NumberOfWaiters; i++ {
		idx := i
		go func() {
			waiters.Waiters[idx].Work(tables, &orderId, rating, conf.KitchenAddress, Menu, &PickUpController)
		}()
	}

	http.ListenAndServe(":"+conf.Port, r)

}

func GetClientOrderDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	clientOrderStatus := clientMap.Get(id, Menu, conf.KitchenAddress)
	resp, _ := json.Marshal(&clientOrderStatus)
	fmt.Fprint(w, string(resp))

}

func RegisterAtOM(rating float64) {
	postBody, _ := json.Marshal(order.RegisterPayload{
		RestaurantId: conf.RestaurantId,
		Address:      conf.LocalAddress,
		Name:         conf.RestaurantName,
		MenuItems:    len(Menu.Foods),
		Menu:         Menu.Foods,
		Rating:       rating,
	})
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://"+conf.OMAddress+"/register", "application/json", responseBody)
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

func GetOrdersStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp, err := http.Get("http://" + conf.KitchenAddress + "/getOrderStatus")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprint(w, string(body))

}

func GetPreparedItems() int {
	resp, err := http.Get("http://" + conf.KitchenAddress + "/getPreparedItems")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	res, _ := strconv.Atoi(string(body))
	return res

}

func ClientRatingPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var ord order.RestaurantRatingPayload
	err := json.NewDecoder(r.Body).Decode(&ord)
	if err != nil {
		log.Fatalln("There was an error decoding the request body into the struct")
	}

	score := rating.AddAndReturn(ord.Rating)

	ratingResponse := order.RestaurantRatingResponse{
		RestaurantId:        conf.RestaurantId,
		RestaurantAvgRating: score,
		PreparedOrders:      GetPreparedItems(),
	}

	resp, _ := json.Marshal(&ratingResponse)
	fmt.Fprint(w, string(resp))
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

	if ord.WaiterId == -1 {
		// Map isering the order
		clientMap.Insert(&ord)
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

	regTime := int(time.Now().UnixMilli())

	var ord = order.Order{
		OrderId:    temp,
		Items:      clientOrd.Items,
		Priority:   clientOrd.Priority,
		MaxWait:    clientOrd.MaxWait,
		PickUpTime: clientOrd.CreatedTime,
		TableId:    regTime,
		WaiterId:   -1,
	}

	// Inserting initiall registration
	go func() {

		clientMap.Insert(&order.Payload{
			OrderId:        ord.OrderId,
			Items:          ord.Items,
			Priority:       ord.Priority,
			MaxWait:        ord.MaxWait,
			TableId:        ord.TableId,
			WaiterId:       ord.WaiterId,
			CookingTime:    0,
			CookingDetails: nil,
		})
	}()

	order.SendOrder(&ord, conf.KitchenAddress)

	clientOrdResponse := order.ClientOrderResponse{
		RestaurantId:         conf.RestaurantId,
		OrderId:              temp,
		EstimatedWaitingTime: order.EstimatedWaitingTimeCalculation(ord, Menu, conf.KitchenAddress),
		CreatedTime:          clientOrd.CreatedTime,
		RegisteredTime:       int64(regTime),
	}

	resp, _ := json.Marshal(&clientOrdResponse)
	fmt.Fprint(w, string(resp))

}
