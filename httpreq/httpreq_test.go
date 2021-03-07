package httpreq

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	ast := assert.New(t)
	html, err := Get("https://example.com/", "")
	ast.NotEmpty(html)
	ast.Nil(err)
}

var resp HTTPResp

func callback(result chan HTTPResp, quit chan int) {
	for {
		select {
		case resp = <-result:
		case <-quit:
			return
		}
	}
}

func TestAsyncGet(t *testing.T) {
	var wg sync.WaitGroup
	ast := assert.New(t)
	result := make(chan HTTPResp)
	quit := make(chan int)
	wg.Add(1)
	go AsyncGet("https://example.com/", "", result, &wg)
	go callback(result, quit)
	wg.Wait()
	quit <- 0
	ast.NotEmpty(resp.body)
	ast.Nil(resp.err)
}
