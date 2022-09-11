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

func main() {
	rand.Seed(time.Now().UnixNano())
	orderId := 0
	sendOrder := func() {

		ord := order.GetRandomOrder(&orderId)
		fmt.Println(orderId)
		time.Sleep(time.Duration(rand.Intn(5-1)+1) * time.Second)
		postBody, _ := json.Marshal(*ord)
		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post("http://localhost:8081/order", "application/json", responseBody)
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

	r := gin.Default()
	r.POST("/distribution", PostHomePage)
	go r.Run(":8080")

	for {
		go sendOrder()
		time.Sleep(time.Duration(rand.Intn(3-1)+1) * time.Second)

	}

}
