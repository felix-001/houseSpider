package main

import (
	"HouseSpider/conf"
	"HouseSpider/request"
	"log"
)

const confPath = "/usr/local/etc/house_spider.conf"

func httpCallback(url, body string, err error, opaque interface{}) {
	log.Println(url)
	//log.Println(err)
	//log.Println(body)
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("enter")
	config, err := conf.Load(confPath)
	if err != nil {
		return
	}
	log.Printf("%+v", config)
	request := request.New(httpCallback)
	request.AsyncGet("http://httpbin.org/get?aaa=123", "", nil)
	request.AsyncGet("http://httpbin.org/get?bbb=456", "", nil)
	request.AsyncGet("http://httpbin.org/get?ccc=789", "", nil)
	request.AsyncGet("http://httpbin.org/get?ccc=111", "", nil)
	request.AsyncGet("http://httpbin.org/get?ccc=222", "", nil)
	request.WaitAllDone()
	//time.Sleep(10 * time.Second)
	//ctx := proxy.New(config.ProxyUrl)
	//ctx.Fetch()
}
