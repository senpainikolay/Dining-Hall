package order

import (
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Table struct {
	TableId         int
	Free            bool
	ReadyToOrder    bool
	WaitingForOrder bool
	Mutex           sync.Mutex
}

type Tables struct {
	Tables []Table
}

func GetTables(nrOfTables int) *Tables {
	var tables Tables

	for i := 1; i <= nrOfTables; i++ {
		tables.Tables = append(tables.Tables,
			Table{
				TableId:         i,
				Free:            true,
				ReadyToOrder:    false,
				WaitingForOrder: false,
				Mutex:           sync.Mutex{},
			})
	}
	return &tables

}

func (ts *Tables) OccupyTables(address string, PickUpControoler *OrderPickUpController) {
	remainder := math.Mod(float64(len(ts.Tables)), float64(len(ts.Tables)/2))

	for i := 0; i < len(ts.Tables); i++ {
		ts.Tables[i].Mutex.Lock()
		if ts.Tables[i].Free && int(remainder) == rand.Intn(len(ts.Tables)/2) {
			ts.Tables[i].Free = false
			tempId := i
			go func() {
				time.Sleep(TIME_UNIT * time.Duration(rand.Intn(5)+5) * time.Millisecond)
				temp := GetOrderStatus(address, PickUpControoler)
				if temp == 0 {
					ts.Tables[tempId].Mutex.Lock()
					ts.Tables[tempId].ReadyToOrder = true
					ts.Tables[tempId].Mutex.Unlock()
				} else {
					ts.Tables[tempId].Mutex.Lock()
					ts.Tables[tempId].Free = true
					ts.Tables[tempId].Mutex.Unlock()

				}
				//log.Printf(" Table %v ready to make the order!", ts.Tables[tempId].TableId)
			}()

			//log.Printf(" Table %v Occupied!", ts.Tables[idx].TableId)
		}
		ts.Tables[i].Mutex.Unlock()
	}
}

func GetOrderStatus(address string, PickUpController *OrderPickUpController) int {
	resp, err := http.Get("http://" + address + "/getOrderStatus")
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	OrderStackIndex, _ := strconv.Atoi(string(body))

	PickUpController.Mutex.Lock()
	PickUpController.SignalVar = OrderStackIndex
	PickUpController.Mutex.Unlock()
	return OrderStackIndex

}
