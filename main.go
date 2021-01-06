package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/corpix/uarand"
)

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.bestbuy.com/site/xfx-amd-radeon-rx-6800xt-16gb-gddr6-pci-express-4-0-gaming-graphics-card-black/6441226.p?skuId=6441226", nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(body))
}
