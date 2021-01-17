package server

import (
	"encoding/json"
	"fmt"
	"goinventory/internal/box"
	"goinventory/models"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"

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
	r.Data = make([]ItemResponse, input.Size)

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

func (self *Server) ServeAbout(w http.ResponseWriter, r *http.Request) {
	aboutTempl := template.Must(template.New("").Parse(string(box.Get("/about.html"))))

	aboutTempl.Execute(w, nil)

}

func (self *Server) ServeFavicon(w http.ResponseWriter, r *http.Request) {
	w.Write(box.Get("/favicon.ico"))

}

func (self *Server) ServeCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.Write(box.Get("/common.css"))

}

func (self *Server) ServeSettings(w http.ResponseWriter, r *http.Request) {

	settingsTemplate := template.Must(template.New("").Parse(string(box.Get("/settings.html"))))

	if r.Method == "GET" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		var settings models.SettingsUpdate = models.SettingsUpdate{
			Delayseconds: self.data.Delayseconds,
			Useragent:    self.data.Useragent,
			Discord:      self.data.Discord,
			Updated:      false,
		}
		settingsTemplate.Execute(w, settings)
	} else if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		r.ParseForm()

		i, err := strconv.Atoi(r.FormValue("irefreshinterval"))
		if err != nil || i < 0 {
			i = int(self.data.Delayseconds)
		}

		var update models.SettingsUpdate = models.SettingsUpdate{
			Delayseconds: int64(i),
			Useragent:    r.FormValue("iuseragent"),
			Discord:      r.FormValue("idiscord"),
			Updated:      true,
		}
		self.data.UpdateFromSettingsUpdate(&update)

		settingsTemplate.Execute(w, update)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
func (self *Server) ServeHome(w http.ResponseWriter, r *http.Request) {
	retVal := InStockResponse{}
	retVal.SetFromSettingsMap(self.data)

	homeTempl := template.Must(template.New("").Parse(string(box.Get("/home.html"))))
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	jsonByte, _ := json.Marshal(retVal.Data)
	var v = struct {
		Data         []ItemResponse
		TimeInterval int
		DataJson     string
	}{Data: retVal.Data, DataJson: string(jsonByte), TimeInterval: int(self.data.Delayseconds)}
	homeTempl.Execute(w, &v)

}

func (self *Server) ServeAddItem(w http.ResponseWriter, r *http.Request) {

	// path, _ := os.Getwd()
	// path = path + "/html/AddItem.html"
	// apage, _ := ioutil.ReadFile(path)
	homeTempl := template.Must(template.New("").Parse(string(box.Get("/AddItem.html"))))
	//homeTempl := template.Must(template.New("").Parse(string(apage)))
	if r.URL.Path != "/add" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if !(r.Method == "GET" || r.Method == "POST") {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()
		_, err := url.ParseRequestURI(r.FormValue("url"))
		if err == nil {
			self.data.AddItem(r.FormValue("iname"), r.FormValue("url"))
			http.Redirect(w, r, "/", http.StatusSeeOther)

		} else {
			fmt.Println("invalid url")
			var errMessage = struct {
				Data string
			}{Data: "Invalid Url"}
			homeTempl.Execute(w, errMessage)
			return
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	homeTempl.Execute(w, nil)
}

func (self *Server) Serve(input *models.SettingsMap, port string) {
	self.data = input
	self.Router = mux.NewRouter().StrictSlash(true)
	self.Router.HandleFunc("/api/items", self.GetInStockItems)
	self.Router.HandleFunc("/", self.ServeHome)
	self.Router.HandleFunc("/favicon.ico", self.ServeFavicon)
	self.Router.HandleFunc("/about", self.ServeAbout)
	self.Router.HandleFunc("/add", self.ServeAddItem)
	self.Router.HandleFunc("/common.css", self.ServeCSS)
	self.Router.HandleFunc("/settings", self.ServeSettings)

	log.Fatal(http.ListenAndServe(":"+port, self.Router))
}
