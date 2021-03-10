package main

import (
	"HouseSpider/conf"
	"HouseSpider/proxy"
	"log"
)

const confPath = "/usr/local/etc/house_spider.conf"

func httpCallback(url, body string, err error, opaque interface{}) {
	log.Println(url)
	log.Println(err)
	log.Println(body)
}

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

	//request := request.New(httpCallback)
	//resp, err := request.Get("https://httpbin.org/get?aaa=bbb", "http://165.227.88.225:8080", nil)
	//request.AsyncGet("https://httpbin.org/get?aaa=bbb", "http://165.227.88.225:8080", nil)
	//request.WaitAllDone()
	//log.Println(err)
	//log.Println(resp)
}
