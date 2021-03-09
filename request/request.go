package request

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const userAget = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"

type Request struct {
	cb HttpCallback
	wg sync.WaitGroup
}

type HttpCallback func(url, body string, err error, opaque interface{})

func New(cb HttpCallback) *Request {
	return &Request{cb: cb, wg: sync.WaitGroup{}}
}

func (r *Request) Get(urlStr, proxyStr string, opaque interface{}) (res string, err error) {
	defer r.wg.Done()
	cli := &http.Client{Timeout: 5 * time.Second}
	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
			return "", err
		}
		transport := &http.Transport{
			Proxy:           http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		cli.Transport = transport
	}
	body := []byte{}
	req, err := http.NewRequest("GET", urlStr, bytes.NewBuffer(body))
	if r.cb != nil {
		defer func() {
			r.cb(urlStr, res, err, opaque)
		}()
	}
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", userAget)
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		err := errors.New(resp.Status)
		return "", err
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println(err)
		return "", err
	}
	return string(body), nil
}

func (r *Request) AsyncGet(urlStr, proxyStr string, opaque interface{}) {
	r.wg.Add(1)
	go r.Get(urlStr, proxyStr, opaque)
}

func (r *Request) WaitAllDone() {
	r.wg.Wait()
}
