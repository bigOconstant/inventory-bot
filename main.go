package main

import (
	"fmt"
	"goinventory/models"
	"goinventory/server"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

func checkStock(wg *sync.WaitGroup, Useragent string, url *models.URLMutex) {
	defer wg.Done()
	client := &http.Client{}

	req, err := http.NewRequest("GET", url.URL, nil)
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
		if !url.InStock {
			url.SetStock(true)
			fmt.Println(url.Name, " in stock go to ", url.URL, " now")
		}

	} else {
		if url.InStock {
			fmt.Println(url.Name, " not in stock now")
		}

	}
}
func main() {
	data := models.SettingsMap{}
	data.ReadFromFile()
	server := server.Server{}
	go server.Serve(&data)
	for true {
		var wg sync.WaitGroup

		//fmt.Println("Spliting into threads")
		for index, item := range data.Items {
			wg.Add(index)
			go checkStock(&wg, data.Useragent, item)
		}

		wg.Wait()
		//fmt.Println("pausing for ", settings.Delayseconds, " seconds")
		time.Sleep(time.Duration(data.Delayseconds) * time.Second)
	}
	fmt.Println("Main: Completed")
}
