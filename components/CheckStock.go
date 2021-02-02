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
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("Error getting page ", url.URL)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Error reading body here")
		println(err.Error())

	}
	// Hack for  stupid best buy comment messing up the algo
	//When you see add to cart, just be patient
	//WHEN YOU SEE ADD TO CART, JUST BE PATIENT
	if strings.Contains(strings.ToUpper((string(body))), "ADD TO CART") && !strings.Contains((string(body)), "When you see add to cart, just be patient") {

		if !url.InStock {
			fmt.Println("in stock")
			fmt.Println(string(body))
			discord.SendNotification(url.Name + " in stock go to " + url.URL + " now")
			url.SetStock(true)
		}

	} else {
		//fmt.Println("not in stock setting")

		if url.InStock {
			fmt.Println("not in stock")
			fmt.Println(url.Name, " not in stock now")
			url.SetStock(false)
		}

	}

	return
}
