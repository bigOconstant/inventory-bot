package main

import (
	"fmt"
	"goinventory/components"
	"goinventory/db/sqlite"
	"goinventory/models"
	"goinventory/server"
	"os"
	"sync"
	"time"
)

func main() {
	var db sqlite.Sqlite
	db.Init()
	defer db.Close()

	port := "3000"
	if len(os.Args) < 2 {
		fmt.Println("Port not specified, defaulting to 3000")
	} else {
		port = os.Args[1]
		fmt.Println("port set to " + port)
	}

	data := models.SettingsMap{}

	data.LoadFromDB(&db)

	discord := components.Discord{Webhook: data.Discord}
	server := server.Server{DB: &db}
	go server.Serve(&data, port)

	//allocate a single waitgroup
	wg := new(sync.WaitGroup)

	// Main loop that checks stock and sleeps a given duration.
	for true {
		if data.Enabled {
			clone := data.Clone()
			for _, item := range clone.Items {
				wg.Add(1)
				go components.CheckStock(wg, clone.Useragent, item, discord)
			}

			wg.Wait()

			time.Sleep(time.Duration(clone.Delayseconds) * time.Second)
			data.Update(clone)
		}

	}
	fmt.Println("Main: Completed")
}
