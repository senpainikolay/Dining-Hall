package order

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Conf struct {
	Port           string `json:"port"`
	KitchenAddress string `json:"kitchen_address"`
}

func GetConf() *Conf {
	jsonFile, err := os.Open("configurations/Conf.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var conf Conf
	json.Unmarshal(byteValue, &conf)
	return &conf

}
