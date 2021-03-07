package main

import (
	"HouseSpider/httpreq"
	"log"
	"sync"
)

var resp httpreq.HTTPResp

func test(result chan httpreq.HTTPResp, quit chan int) {
	for {
		select {
		case r := <-result:
			log.Println(r)
		case <-quit:
			log.Println("quit")
			return
		}
	}
	resp = <-result
	log.Println(resp)
}

func main() {
	log.SetFlags(log.Lshortfile)
	log.Println("enter")
	result := make(chan httpreq.HTTPResp)
	quit := make(chan int)
	var wg sync.WaitGroup
	go test(result, quit)
	wg.Add(1)
	go httpreq.AsyncGet("https://example.com/", "", result, &wg)
	wg.Wait()
	quit <- 0
	log.Println("done")

}
