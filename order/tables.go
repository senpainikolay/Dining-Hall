package order

import (
	"log"
	"math"
	"math/rand"
)

type Table struct {
	TableId         int
	Free            bool
	ReadToOrder     bool
	WaitingForOrder bool
	OrderRecieving  chan Order
}

func OccupyTables(tables []Table, nrOfTables int) {
	remainder := math.Mod(float64(nrOfTables), float64(nrOfTables/2))
	for i := 0; i < nrOfTables; i++ {
		if tables[i].Free == true && int(remainder) == rand.Intn(nrOfTables/2) {
			tables[i].Free = false
			log.Printf(" Table %v Occupied!", tables[i].TableId)
		}
	}

}
