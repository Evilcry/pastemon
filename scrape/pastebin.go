package scrape

import (
	"encoding/json"
	"log"
	"time"

	"github.com/evilcry/pastemon/netutils"
	"github.com/evilcry/pastemon/processing"

	"github.com/evilcry/pastemon/configs"
)

type Paste struct {
	ScrapeUrl string `json:"scrape_url"`
	Url       string `json:"full_url"`
	Date      string
	Key       string
	Size      int `json:",string"`
	Expire    int `json:",string"`
	Title     string
	Syntax    string
	User      string
	Error     string
	Content   string
}

// Download func
// retrieve paste
func (p *Paste) Download(conf *configs.Config) {
	_, exists := conf.Keys[p.Key]
	if exists {
		return
	}

	resp := netutils.Get(p.ScrapeUrl)
	p.Content = string(resp)
	conf.Keys[p.Key] = time.Now()
}

// Process func
// processes paste content
func (p *Paste) Process(conf *configs.Config) {
	processing.ProcessContent(conf, p.Key, p.Content)
}

// PastebinScraper func
// scraper
func PastebinScraper(conf *configs.Config) {
	var pastes []*Paste

	log.Println("[+] Checking for new pastes.")

	resp := netutils.Get("https://scrape.pastebin.com/api_scraping.php?limit=100")
	err := json.Unmarshal(resp, &pastes)
	if err != nil {
		log.Println("[-] Could not parse list of pastes.")
		log.Printf("[-] %s.\n", err.Error())
		log.Println(string(resp))
		return
	}

	for i, _ := range pastes {
		p := pastes[i]
		p.Download(conf)
		p.Process(conf)
	}
}
