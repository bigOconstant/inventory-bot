package components

import (
	"fmt"
	"goinventory/models"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

func CheckStock(wg *sync.WaitGroup, Useragent string, url *models.URLMutex, discord Discord) {
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
			discord.SendNotification(url.Name + " in stock go to " + url.URL + " now")
		}

	} else {
		if url.InStock {
			url.SetStock(false)
			fmt.Println(url.Name, " not in stock now")
		}

	}
}
