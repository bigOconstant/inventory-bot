package server

import (
	"encoding/json"
	"fmt"
	"goinventory/db/dbmodels"
	"goinventory/db/sqlite"
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
	input.Lock()
	r.Data = make([]ItemResponse, len(input.Items))

	i := 0
	for key := range input.Items {
		r.Data[i].Id = i
		r.Data[i].Name = input.Items[key].Name
		r.Data[i].Url = input.Items[key].URL
		r.Data[i].InStock = input.Items[key].InStock
		i++
	}
	input.Unlock()
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
			Enabled:      self.data.Enabled,
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
		enabled := r.FormValue("enabled")
		enabledBool := false
		if enabled == "on" {
			enabledBool = true
		}

		var update models.SettingsUpdate = models.SettingsUpdate{
			Delayseconds: int64(i),
			Useragent:    r.FormValue("iuseragent"),
			Discord:      r.FormValue("idiscord"),
			Enabled:      enabledBool,
			Updated:      true,
		}
		self.data.UpdateFromSettingsUpdate(&update)

		settingsUpdate := dbmodels.Settings{Id: 1,
			Refresh_interval: int(update.Delayseconds),
			User_agent:       update.Useragent,
			Enabled:          update.Enabled,
			Discord_webhook:  update.Discord}

		fmt.Println(settingsUpdate)
		db := sqlite.Sqlite{}
		db.SaveSettings(settingsUpdate)

		settingsTemplate.Execute(w, update)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func (self *Server) ServeHome(w http.ResponseWriter, r *http.Request) {

	homeTempl := template.Must(template.New("").Parse(string(box.Get("/home.html"))))
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == "POST" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		r.ParseForm()
		for _, values := range r.Form { // range over map
			for _, value := range values { // range over []string
				id, err := strconv.Atoi(value)
				if err == nil {
					self.data.RemoveID(id)
				}
			}
		}
	}
	retVal := InStockResponse{}
	retVal.SetFromSettingsMap(self.data)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	jsonByte, _ := json.Marshal(retVal.Data)
	var v = struct {
		Data         []ItemResponse
		TimeInterval int
		DataJson     string
		Enabled      bool
	}{Data: retVal.Data, DataJson: string(jsonByte), TimeInterval: int(self.data.Delayseconds), Enabled: self.data.Enabled}
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
			db := sqlite.Sqlite{}
			id, err := db.SaveItem(r.FormValue("iname"), r.FormValue("url"))
			fmt.Println("new id :", id, " err:", err)
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
