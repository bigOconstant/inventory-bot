package models

import (
	"encoding/json"
	"fmt"
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
	Size         int
	Enabled      bool
	Items        map[int]*URLMutex
}

func (s *SettingsMap) Clone() *SettingsMap {
	s.mu.Lock()
	defer s.mu.Unlock()
	retVal := SettingsMap{Size: s.Size, Delayseconds: s.Delayseconds, Useragent: s.Useragent, Discord: s.Discord, Enabled: s.Enabled}

	retVal.Items = make(map[int]*URLMutex, s.Size)
	for i := 0; i < s.Size; i++ {
		retVal.Items[i] = &URLMutex{URL: s.Items[i].URL, Name: s.Items[i].Name, Id: s.Items[i].Id, InStock: s.Items[i].InStock}
	}
	return &retVal
}

func (old *SettingsMap) Update(new *SettingsMap) {
	for _, newurl := range new.Items {
		for _, oldurl := range old.Items {
			if newurl.Name == oldurl.Name && newurl.URL == oldurl.URL {
				oldurl.InStock = newurl.InStock
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
	s.Size = len(input.Urls)
	s.Useragent = input.Useragent
	s.Items = make(map[int]*URLMutex)
	s.Discord = input.Discord
	s.Enabled = true
	for i := 0; i < s.Size; i++ {
		s.Items[i] = &URLMutex{Id: i}
		s.Items[i].SetFromUrls(input.Urls[i])
	}
}

func (s *SettingsMap) AddItem(name string, url string) {
	length := len(s.Items)
	var UrlModel Urls = Urls{Item: name, URL: url}
	urlMutex := URLMutex{}
	urlMutex.SetFromUrls(UrlModel)
	s.Items[length] = &urlMutex
	s.Items[length].mu.Lock()
	defer s.Items[length].mu.Unlock()
	s.Items[length].URL = url
	s.Items[length].Name = name
	s.Items[length].Id = length
	s.Size++
}

func (u *SettingsMap) ReadFromFile() {
	settingsFile, err := os.Open("settings.json")
	settings := Settings{}
	if err != nil {
		fmt.Println(err)
	}

	byteValue, _ := ioutil.ReadAll(settingsFile)
	json.Unmarshal([]byte(byteValue), &settings)
	settingsFile.Close()
	u.FromSettings(&settings)
}

func (u *SettingsMap) Lock() {
	u.mu.Lock()
}

func (u *SettingsMap) Unlock() {
	u.mu.Unlock()
}
