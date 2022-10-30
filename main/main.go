package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"sync"
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

var PickUpController = order.OrderPickUpController{Mutex: sync.Mutex{}, SignalVar: 0}

func main() {

	rand.Seed(time.Now().UnixMilli())
	tables = order.GetTables(NumberOfTables)
	waiters = order.GetWaiters(NumberOfWaiters)
	rating := order.GetRating()

	r := mux.NewRouter()
	r.HandleFunc("/distribution", PostKitchenOrders).Methods("POST")
	r.HandleFunc("/v2/order", PostClientOrders).Methods("POST")
	r.HandleFunc("/getOrderStatus", GetOrdersStatus).Methods("GET")

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
		EstimatedWaitingTime: EstimatedWaitingTimeCalculation(ord),
		CreatedTime:          clientOrd.CreatedTime,
		RegisteredTime:       time.Now().UnixMilli(),
	}

	resp, _ := json.Marshal(&clientOrdResponse)
	fmt.Fprint(w, string(resp))

}

func EstimatedWaitingTimeCalculation(ord order.Order) float64 {

	var wg sync.WaitGroup
	wg.Add(2)
	var A, B, C, D, E, F float64
	go func() {
		for _, food_id := range ord.Items {
			if Menu.Foods[food_id-1].CookingApparatus == "" {
				A += float64(Menu.Foods[food_id-1].PreparationTime)
				continue
			}
			C += float64(Menu.Foods[food_id-1].PreparationTime)
		}

		F = float64(len(ord.Items))
		wg.Done()
	}()

	go func() {
		resp, err := http.Get("http://" + conf.KitchenAddress + "/estimationCalculation")
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var BDE order.EWTCalculation
		if err := json.Unmarshal([]byte(body), &BDE); err != nil {
			panic(err)
		}
		B = float64(BDE.B)
		D = float64(BDE.D)
		E = float64(BDE.E)
		wg.Done()
	}()

	wg.Wait()
	fm := ((A/B + C/D) * (E + F)) / F

	return roundFloat(fm, 2)

	// log.Println(menuInfo)

}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
