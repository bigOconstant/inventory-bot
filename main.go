package main

import (
	"fmt"
	"goinventory/models"
	"goinventory/server"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func checkStock(wg *sync.WaitGroup, Useragent string, url *models.URLMutex) {
	defer wg.Done()
	client := &http.Client{}

	req, err := http.NewRequest("GET", url.URL, nil)
	if err != nil {
		println("Error calling Get")
		println(err.Error())
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", Useragent)

	resp, err := client.Do(req)
	if err != nil {
		println("Error calling do!")
		println(err.Error())
		log.Fatalln(err)

	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error getting page ", url.URL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error reading body here")
		println(err.Error())
		log.Fatalln(err)
	}

	if strings.Contains(strings.ToUpper((string(body))), "ADD TO CART") {
		if !url.InStock {
			url.SetStock(true)
			fmt.Println(url.Name, " in stock go to ", url.URL, " now")
		}

	} else {
		if url.InStock {
			url.SetStock(false)
			fmt.Println(url.Name, " not in stock now")
		}

	}
}
func main() {
	port := "3000"
	if len(os.Args) < 2 {
		fmt.Println("Port not specified, defaulting to 3000")
	} else {
		port = os.Args[1]
		fmt.Println("port set to " + port)
	}

	data := models.SettingsMap{}
	data.ReadFromFile()
	server := server.Server{}
	go server.Serve(&data, port)
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
