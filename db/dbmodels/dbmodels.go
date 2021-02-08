package dbmodels

type Settings struct {
	Id               int
	Refresh_interval int
	User_agent       string
	Discord_webhook  string
	Enabled          bool
}

type Item struct {
	Id   int
	Url  string
	Name string
}
