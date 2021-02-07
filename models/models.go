package models

import (
	"encoding/json"
	"fmt"
	"goinventory/db/sqlite"
	"io/ioutil"
	"os"
	"sync"
)

/*
Urls Holds item to pass to each worker
*/
type Urls struct {
	Item string `json:"item"`
	URL  string `json:"url"`
}

type DiscordMessage struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

/*
Settings holds settings json file
*/
type Settings struct {
	Delayseconds int64  `json:"delayseconds"`
	Useragent    string `json:"useragent"`
	Urls         []Urls `json:"urls"`
	Discord      string `json:"discord"`
}

type SettingsUpdate struct {
	Delayseconds int64
	Useragent    string
	Discord      string
	Updated      bool
	Enabled      bool
}

//Url struct
type URLMutex struct {
	mu      sync.Mutex
	URL     string
	Id      int
	InStock bool
	Name    string
}

//SetStock sets stock thread safe
func (u *URLMutex) SetStock(input bool) {
	u.mu.Lock()
	u.InStock = input
	u.mu.Unlock()
}

//SetFromUrls sets mutex struct from a url struct
func (u *URLMutex) SetFromUrls(input Urls) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.InStock = false
	u.URL = input.URL
	u.Name = input.Item
}

type SettingsMap struct {
	mu           sync.Mutex
	Delayseconds int64
	Useragent    string
	Discord      string
	Enabled      bool
	Items        map[int]*URLMutex
}

//RemoveID removes an item from the map
func (s *SettingsMap) RemoveID(id int) {

	newItems := make(map[int]*URLMutex, len(s.Items))

	iterator := 0

	for i, item := range s.Items {
		if i != id {
			newItems[iterator] = item
			iterator++
		}
	}
	s.Items = nil
	s.Items = newItems

}

func (s *SettingsMap) Clone() *SettingsMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	retVal := SettingsMap{Delayseconds: s.Delayseconds, Useragent: s.Useragent, Discord: s.Discord, Enabled: s.Enabled}

	retVal.Items = make(map[int]*URLMutex, len(s.Items))
	for key, _ := range s.Items {
		retVal.Items[key] = &URLMutex{URL: s.Items[key].URL, Name: s.Items[key].Name, Id: s.Items[key].Id, InStock: s.Items[key].InStock}
	}
	return &retVal
}

func (old *SettingsMap) Update(new *SettingsMap) {

	for _, newurl := range new.Items {
		for _, oldurl := range old.Items {
			if newurl.Name == oldurl.Name && newurl.URL == oldurl.URL {
				//oldurl.mu.Lock()
				oldurl.InStock = newurl.InStock
				//oldurl.mu.Unlock()
			}
		}
	}
}

func (s *SettingsMap) UpdateFromSettingsUpdate(su *SettingsUpdate) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Delayseconds = su.Delayseconds
	s.Discord = su.Discord
	s.Useragent = su.Useragent
	s.Enabled = su.Enabled
}

func (s *SettingsMap) FromSettings(input *Settings) {
	s.Delayseconds = input.Delayseconds
	s.Useragent = input.Useragent
	s.Items = make(map[int]*URLMutex)
	s.Discord = input.Discord
	s.Enabled = true
	for i := 0; i < len(input.Urls); i++ {
		s.Items[i] = &URLMutex{Id: i}
		s.Items[i].SetFromUrls(input.Urls[i])
	}
}

func (s *SettingsMap) AddItem(name string, url string) {

	newItems := make(map[int]*URLMutex, len(s.Items)+1)

	i := 0
	for _, item := range s.Items {
		newItems[i] = item
		i++
	}

	var UrlModel Urls = Urls{Item: name, URL: url}

	um := URLMutex{}
	um.SetFromUrls(UrlModel)
	s.Items = nil
	newItems[i] = &um
	s.Items = newItems

}

func (u *SettingsMap) ReadFromFile() {
	settingsFile, err := os.Open("settings.json")
	settings := Settings{}
	if err != nil {
		settings.Useragent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
		settings.Delayseconds = 30
		u.FromSettings(&settings)
		fmt.Println(err)
		return
	}

	byteValue, _ := ioutil.ReadAll(settingsFile)
	json.Unmarshal([]byte(byteValue), &settings)
	settingsFile.Close()
	u.FromSettings(&settings)
}

func (U *SettingsMap) LoadFromDB(db *sqlite.Sqlite) {
	settings, _ := db.GetSettings()
	U.Delayseconds = int64(settings.Refresh_interval)
	U.Discord = settings.Discord_webhook
	U.Enabled = settings.Enabled
	U.Useragent = settings.User_agent
	U.Items = nil
	U.Items = make(map[int]*URLMutex)
	items, err := db.GetItems()
	if err == nil {
		for _, item := range items {
			var u URLMutex = URLMutex{URL: item.Url, Name: item.Name, InStock: false}
			U.Items[item.Id] = &u
		}
	}
}

func (u *SettingsMap) Lock() {
	u.mu.Lock()
}

func (u *SettingsMap) Unlock() {
	u.mu.Unlock()
}
