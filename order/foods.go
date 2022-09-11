package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Foods struct {
	Foods []Food `json:"foods"`
}

type Food struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	PreparationTime  int    `json:"preparation_time"`
	Complexity       int    `json:"complexity"`
	CookingApparatus string `json:"cooking_apparatus"`
}

func getFoods() *Foods {
	jsonFile, err := os.Open("configurations/Foods.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var menu Foods
	json.Unmarshal(byteValue, &menu)
	return &menu

}
