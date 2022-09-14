package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	dbAddr = "www.douban.com/group/search"
)

//var groups = []string{"26926", "279962", "262626", "35417", "56297", "257523", "374051", "625354", "aihezu", "zhufang", "opking", "jumei", "beijingzufang"}
var groups = []string{"26926"}
var keywords = []string{"两居"}

func httpGet(url string) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("http resp status %s code: %d ", resp.Status, resp.StatusCode)
	}
	return string(body), nil
}

type Spider struct {
}

func New() *Spider {
	return &Spider{}
}

func (s *Spider) Run() error {
	for _, keyword := range keywords {
		for _, group := range groups {
			searchAddr := fmt.Sprintf("https://%s?group=%s&cat=1013&q=%s", dbAddr, group, keyword)
			resp, err := httpGet(searchAddr)
			if err != nil {
				return err
			}
			if err := s.parse(resp); err != nil {
				return err
			}
		}
	}
	return nil
}

func callback(i int, s *goquery.Selection) {
	title := s.Find("a").Text()
	fmt.Printf("Review %d: %s\n", i, title)
}

func (s *Spider) parse(html string) error {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBufferString(html))
	if err != nil {
		return err
	}

	doc.Find(".pl").Each(callback)
	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	s := New()
	if err := s.Run(); err != nil {
		log.Println(err)
	}
}
