package main

import (
	"fmt"
	"goinventory/components"
	"goinventory/models"
	"goinventory/server"
	"os"
	"sync"
	"time"
)

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
	discord := components.Discord{Webhook: data.Discord}
	server := server.Server{}
	go server.Serve(&data, port)

	// Main loop that checks stock and sleeps a given duration.
	for true {
		var wg sync.WaitGroup
		for index, item := range data.Items {
			wg.Add(index)
			go components.CheckStock(&wg, data.Useragent, item, discord)
		}
		wg.Wait()
		time.Sleep(time.Duration(data.Delayseconds) * time.Second)
	}
	fmt.Println("Main: Completed")
}
