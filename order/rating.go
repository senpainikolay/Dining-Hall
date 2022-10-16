package order

import (
	"log"
	"sync"
)

type Rating struct {
	Sum     float64
	Counter float64
	Score   float64
	sync.Mutex
}

func GetRating() *Rating {
	return &Rating{0.0, 0.0, 0.0, sync.Mutex{}}
}

func (r *Rating) Reformulate(maxWait float64, timeServed float64) int {

	if timeServed <= maxWait {
		return 5
	}
	if timeServed <= maxWait*1.1 {
		return 4
	}
	if timeServed <= maxWait*1.2 {
		return 3
	}
	if timeServed <= maxWait*1.3 {
		return 2
	}
	if timeServed <= maxWait*1.4 {
		return 1
	}
	return 0

}

func (r *Rating) Calculate(maxWait float64, timeServer float64, address string) {
	nr := r.Reformulate(maxWait, timeServer)
	r.Sum += float64(nr)
	r.Counter += 1
	r.Score = r.Sum / r.Counter
	log.Printf("The actual score: %v at %v \n", r.Score, address)
}
