package order

import (
	"log"
	"math"
	"math/rand"
	"sync"
	"time"
)

type Table struct {
	TableId         int
	Free            bool
	ReadyToOrder    bool
	WaitingForOrder bool
}

type Tables struct {
	Tables []Table
	Mutex  sync.Mutex
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
			})
	}
	return &tables

}

func (ts *Tables) OccupyTables() {
	remainder := math.Mod(float64(len(ts.Tables)), float64(len(ts.Tables)/2))
	ts.Mutex.Lock()
	defer ts.Mutex.Unlock()

	for i := 0; i < len(ts.Tables); i++ {
		idx := i
		go func() {
			if ts.Tables[idx].Free && int(remainder) == rand.Intn(len(ts.Tables)/2) {
				ts.Tables[idx].Free = false
				tempId := idx
				go func() {
					time.Sleep(2 * time.Second)
					ts.Tables[tempId].ReadyToOrder = true
					log.Printf(" Table %v ready to make the order!", ts.Tables[tempId].TableId)
				}()
				log.Printf(" Table %v Occupied!", ts.Tables[idx].TableId)
			}
		}()
	}
}
