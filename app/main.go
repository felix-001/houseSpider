package main

import (
	"HouseSpider/conf"
	"HouseSpider/house"
	"HouseSpider/proxy"
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
	proxy := proxy.New(config.ProxyUrl)
	proxy.Fetch()
	proxy.Filter()
	house := house.New(config, proxy)
	house.Fetch()
}
