package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/evilcry/pastemon/configs"
	"github.com/evilcry/pastemon/db"
	"github.com/evilcry/pastemon/scrape"
)

func main() {
	log.Println("--=Pastebin Monitoring Service=--")

	conf := configs.NewConfig()

	ds := db.GetStorageConnection(&conf)
	db.InitStorage(&conf, ds)
	ds.Close()

	log.Println("[+] Pastemon is scraping...")

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("[-] Shutting Down Pastemon")
		os.Exit(0)
	}()

	for {
		scrape.PastebinScraper(&conf)
		time.Sleep(time.Duration(conf.Sleep) * time.Second)
		db.CleanKeys(&conf)
	}

}
