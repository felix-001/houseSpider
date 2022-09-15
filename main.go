package main

import (
	"HouseSpider/house"
	"log"
)

const confPath = "/usr/local/etc/house_spider.conf"

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("enter")
	/*
		config, err := conf.Load(confPath)
		if err != nil {
			return
		}
		log.Printf("%+v", config)
		h := house.New(config)
		h.Fetch()
	*/
	proxyMgr := house.NewProxyManager()
	proxyMgr.Register(&house.SPYSAgent{})
	proxyMgr.Init()
	proxyMgr.Run()
}
