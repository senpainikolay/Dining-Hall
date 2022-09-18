package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/senpainikolay/Dining-Hall/order"
)

func PostHomePage(c *gin.Context) {

	body := c.Request.Body
	value, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(200, gin.H{
		"message": string(value),
	})

	fmt.Printf("Post success: %s \n", value)
}

const (
	NumberOfTables  = 10
	NumberOfWaiters = 5
)

func main() {
	rand.Seed(time.Now().UnixNano())
	/*
		orderId := 0

			sendOrder := func() {

				ord := order.GetRandomOrder(&orderId)
				fmt.Println(orderId)
				time.Sleep(time.Duration(rand.Intn(5-1)+1) * time.Second)
				postBody, _ := json.Marshal(*ord)
				responseBody := bytes.NewBuffer(postBody)
				resp, err := http.Post("http://kitchen:8081/order", "application/json", responseBody)
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
	r := gin.Default()
	r.POST("/distribution", PostHomePage)
	go r.Run(":8080")

	var tables []order.Table

	for i := 1; i <= NumberOfTables; i++ {
		tables = append(tables,
			order.Table{
				TableId:         i,
				Free:            true,
				ReadyToOrder:    false,
				WaitingForOrder: false,
				OrderRecieving:  make(chan order.Order, 1),
			})
	}

	var waiters []order.Waiter

	for i := 1; i <= NumberOfWaiters; i++ {
		waiters = append(waiters,
			order.Waiter{
				WaiterId:        i,
				OrdersToRecieve: make(chan order.Order),
				OrdersToServe:   make(chan order.Order),
			})
	}

	order.OccupyTables(tables, NumberOfTables)
	log.Println("AUFFF")
	orderId := order.OrderId{Id: 0}

	for {

		for _, waiter := range waiters {

			select {
			case PostOrder := <-waiter.OrdersToRecieve:
				fmt.Printf("%+v To send to Kitchen from waiter Id %v", PostOrder, waiter.WaiterId)

			case ServeOrder := <-waiter.OrdersToServe:
				fmt.Printf("%+v Serving", ServeOrder)

			default:
				time.Sleep(100 * time.Millisecond)
				waiter.PickUpOrder(tables, &orderId)

			}

		}

	}

	// ADDED FOR MAIN

	// COMMENT FOR TEST BRANCH
	/*
		for {
			go sendOrder()
			time.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

		}

	*/

}
