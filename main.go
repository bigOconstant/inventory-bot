package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

/*
Urls Holds item to pass to each worker
*/
type Urls struct {
	Item string `json:"item"`
	URL  string `json:"url"`
}

/*
Settings holds settings json file
*/
type Settings struct {
	Delayseconds int64  `json:"delayseconds"`
	Useragent    string `json:"useragent"`
	Urls         []Urls `json:"urls"`
}

func main() {
	settingsFile, err := os.Open("settings.json")
	settings := Settings{}
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(settingsFile)
	json.Unmarshal([]byte(byteValue), &settings)
	settingsFile.Close()

	for true {
		var wg sync.WaitGroup
		//fmt.Println("Spliting into threads")
		for index, item := range settings.Urls {
			wg.Add(index)
			go checkStock(&wg, index, item, settings.Useragent)
		}

		//fmt.Println("Waiting for workers to finish")
		wg.Wait()

		//fmt.Println("pausing for ", settings.Delayseconds, " seconds")
		time.Sleep(time.Duration(settings.Delayseconds) * time.Second)
	}
	fmt.Println("Main: Completed")
}

func checkStock(wg *sync.WaitGroup, id int, item Urls, Useragent string) {
	defer wg.Done()
	client := &http.Client{}

	req, err := http.NewRequest("GET", item.URL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", Useragent)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if strings.Contains(strings.ToUpper((string(body))), "ADD TO CART") {
		fmt.Println(item.Item, " in stock go to ", item.URL, " now")

	} else {
		fmt.Println(item.Item, " not in stock")
	}
}
