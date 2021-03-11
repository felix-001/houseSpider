package main

import (
	"HouseSpider/conf"
	"HouseSpider/house"
	"log"
)

const confPath = "/usr/local/etc/house_spider.conf"

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("enter")
	config, err := conf.Load(confPath)
	if err != nil {
		return
	}
	log.Printf("%+v", config)
	h := house.New(config)
	h.Fetch()
}
