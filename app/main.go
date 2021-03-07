package main

import (
	"HouseSpider/conf"
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
	proxyCtx := proxy.New(config.ProxyUrl)
	proxyCtx.Fetch()
	proxy, _ := proxyCtx.Get()
	log.Println(proxy)
}
