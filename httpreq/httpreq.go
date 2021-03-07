package httpreq

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

const userAget = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"

func Get(urlStr, proxyStr string) (string, error) {
	cli := &http.Client{}
	if proxyStr != "" {
		proxyURL, err := url.Parse(proxyStr)
		if err != nil {
			log.Println(err)
			return "", err
		}
		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		cli.Transport = transport
	}
	body := []byte{}
	req, err := http.NewRequest("GET", urlStr, bytes.NewBuffer(body))
	if err != nil {
		log.Println(err)
		return "", err
	}
	req.Header.Add("User-Agent", userAget)
	resp, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("status code error: %d %s", resp.StatusCode, resp.Status)
		return "", errors.New(resp.Status)
	}
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println(err)
		return "", err
	}
	return string(body), nil
}

type HTTPResp struct {
	body string
	err  error
}

func AsyncGet(urlStr, proxyStr string, result chan HTTPResp, wg *sync.WaitGroup) {
	body, err := Get(urlStr, proxyStr)
	resp := HTTPResp{body: body, err: err}
	result <- resp
	wg.Done()
}
