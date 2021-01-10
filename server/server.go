package server

import (
	"encoding/json"
	"fmt"
	"goinventory/models"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type ItemResponse struct {
	Url     string `json:"url"`
	Id      int    `json:id`
	Name    string `json:"name"`
	InStock bool   `json:"instock"`
}

type InStockResponse struct {
	Data []ItemResponse `json:"data"`
}

func (r *InStockResponse) SetFromSettingsMap(input *models.SettingsMap) {
	r.Data = make([]ItemResponse, len(input.Items))

	for i := 0; i < input.Size; i++ {
		r.Data[i].Id = i
		r.Data[i].Name = input.Items[i].Name
		r.Data[i].Url = input.Items[i].URL
		r.Data[i].InStock = input.Items[i].InStock
	}
}

type Server struct {
	Router *mux.Router
	data   *models.SettingsMap
}

func (self *Server) GetInStockItems(w http.ResponseWriter, r *http.Request) {
	retVal := InStockResponse{}
	retVal.SetFromSettingsMap(self.data)
	msg, _ := json.Marshal(retVal)
	fmt.Fprintf(w, string(msg))

}
func (self *Server) faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/favicon.ico")
}

func (self *Server) logFileHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./logfile")
}

func (self *Server) ServeTestFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/test.html")
}
func (self *Server) ServeAbout(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/about.html")
}

func (self *Server) ServeHome(w http.ResponseWriter, r *http.Request) {
	path, _ := os.Getwd()
	path = path + "/html/index.html"
	hpage, _ := ioutil.ReadFile(path)
	homeTempl := template.Must(template.New("").Parse(string(hpage)))
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var v = struct {
		Host string
		Port string
	}{Host: self.data.Host, Port: self.data.Port}
	homeTempl.Execute(w, &v)
}

func (self *Server) Serve(input *models.SettingsMap, port string) {
	self.data = input
	self.Router = mux.NewRouter().StrictSlash(true)
	self.Router.HandleFunc("/api/items", self.GetInStockItems)
	self.Router.HandleFunc("/", self.ServeHome)
	self.Router.HandleFunc("/logs", self.logFileHandler)
	self.Router.HandleFunc("/test", self.ServeTestFile)
	self.Router.HandleFunc("/about", self.ServeAbout)
	self.Router.HandleFunc("/favicon.ico", self.faviconHandler)

	log.Fatal(http.ListenAndServe(":"+port, self.Router))
}
