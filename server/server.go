package server

import (
	"encoding/json"
	"fmt"
	"goinventory/models"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type ItemResponse struct {
	Url  string `json:"url"`
	Id   int    `json:id`
	Name string `json:"name"`
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

func (self *Server) Serve(input *models.SettingsMap) {
	self.data = input
	self.Router = mux.NewRouter().StrictSlash(true)
	self.Router.HandleFunc("/items", self.GetInStockItems)
	log.Fatal(http.ListenAndServe(":3000", self.Router))
}
