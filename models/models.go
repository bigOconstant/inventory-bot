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
	Host         string `json:"host"`
	Port         string `json:"port"`
	Discord      string `json:"discord"`
}

//Url struct
type URLMutex struct {
	mu      sync.Mutex
	URL     string
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
	Delayseconds int64
	Useragent    string
	Host         string
	Port         string
	Discord      string
	Size         int
	Items        map[int]*URLMutex
}

func (s *SettingsMap) FromSettings(input *Settings) {
	s.Delayseconds = input.Delayseconds
	s.Size = len(input.Urls)
	s.Useragent = input.Useragent
	s.Items = make(map[int]*URLMutex)
	s.Host = input.Host
	s.Port = input.Port
	s.Discord = input.Discord
	for i := 0; i < s.Size; i++ {
		s.Items[i] = &URLMutex{}
		s.Items[i].SetFromUrls(input.Urls[i])
	}
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
