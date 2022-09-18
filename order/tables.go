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
	OrderRecieving  chan Order
	Mutex           sync.Mutex
}

func OccupyTables(tables []Table, nrOfTables int) {
	remainder := math.Mod(float64(nrOfTables), float64(nrOfTables/2))

	for i := 0; i < nrOfTables; i++ {
		idx := i
		tables[idx].Mutex.Lock()
		defer tables[idx].Mutex.Unlock()
		if tables[idx].Free == true && int(remainder) == rand.Intn(nrOfTables/2) {
			tables[idx].Free = false
			go func() {
				time.Sleep(3 * time.Second)
				tables[idx].ReadyToOrder = true
				log.Printf(" Table %v ready to make the order!", tables[idx].TableId)
			}()
			log.Printf(" Table %v Occupied!", tables[idx].TableId)
		}

	}
}
