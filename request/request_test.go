package request

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func callback(body string, err error) {
	log.Println(body)
	log.Println(err)
}

func TestGet(t *testing.T) {
	ast := assert.New(t)
	request := NewRequest(callback)
	html, err := request.Get("http://httpbin.org/get", "")
	ast.NotEmpty(html)
	ast.Nil(err)
}
