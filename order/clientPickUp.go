package order

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"sync"
)

type GetClientOrder struct {
	OrderMap map[int]Payload
	Mutex    sync.Mutex
}

func GetClientMap() *GetClientOrder {
	return &GetClientOrder{
		OrderMap: make(map[int]Payload),
		Mutex:    sync.Mutex{},
	}
}

func (cm *GetClientOrder) Insert(py *Payload) {
	cm.Mutex.Lock()
	cm.OrderMap[py.OrderId] = *py
	cm.Mutex.Unlock()
}

func (cm *GetClientOrder) Get(OrderId int, Menu *Foods, kitchenAddr string) *ClientOrderStatus {
	cm.Mutex.Lock()

	ord := cm.OrderMap[OrderId]
	// in case it is cooked
	if ord.CookingDetails != nil {
		cm.Mutex.Unlock()
		delete(cm.OrderMap, OrderId)
		return &ClientOrderStatus{
			OrderId:              ord.OrderId,
			IsReady:              true,
			EstimatedWaitingTime: 0.0,
			Priority:             ord.Priority,
			MaxWait:              ord.MaxWait,
			CreatedTime:          ord.PickUpTime,
			RegisteredTime:       ord.TableId, // just couldnt store it somewhere else :-)
			PreparedTime:         ord.CookingTime,
			CookingTime:          ord.CookingTime,
			CookingDetails:       ord.CookingDetails,
		}
	}
	cm.Mutex.Unlock()

	orderValidTypeForEWT := Order{
		OrderId:    ord.OrderId,
		Items:      ord.Items,
		Priority:   ord.Priority,
		MaxWait:    ord.MaxWait,
		PickUpTime: ord.PickUpTime,
		TableId:    ord.TableId,
		WaiterId:   ord.WaiterId,
	}
	return &ClientOrderStatus{
		OrderId:              ord.OrderId,
		IsReady:              false,
		EstimatedWaitingTime: EstimatedWaitingTimeCalculation(orderValidTypeForEWT, Menu, kitchenAddr),
		Priority:             ord.Priority,
		MaxWait:              ord.MaxWait,
		CreatedTime:          ord.PickUpTime,
		RegisteredTime:       ord.TableId, // just couldnt store it somewhere else :-)
		PreparedTime:         0,
		CookingTime:          0,
		CookingDetails:       nil,
	}

}

type ClientOrderStatus struct {
	OrderId              int              `json:"order_id"`
	IsReady              bool             `json:"is_ready"`
	EstimatedWaitingTime float64          `json:"estimated_waiting_time"`
	Priority             int              `json:"priority"`
	MaxWait              float64          `json:"max_wait"`
	CreatedTime          int64            `json:"created_time"`
	RegisteredTime       int              `json:"registered_time"`
	PreparedTime         int64            `json:"prepared_time"`
	CookingTime          int64            `json:"cooking_time"`
	CookingDetails       []CookingDetails `json:"cooking_details"`
}

func EstimatedWaitingTimeCalculation(ord Order, Menu *Foods, kitchenAddr string) float64 {

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
		resp, err := http.Get("http://" + kitchenAddr + "/estimationCalculation")
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}

		var BDE EWTCalculation
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

}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
