package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	dbAddr    = "www.douban.com/group/search"
	pageCount = 5
	sleepTime = 5
)

//var groups = []string{"26926", "279962", "262626", "35417", "56297", "257523", "374051", "625354", "aihezu", "zhufang", "opking", "jumei", "beijingzufang"}
var groups = []string{"26926"}
var keywords = []string{"两居"}
var filterKeywords = []string{"限女生", "求租", "限女", "合租", "室友"}

var result string

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
			for i := 0; i < pageCount; i++ {
				searchAddr := fmt.Sprintf("https://%s?group=%s&cat=1013&q=%s&start=%d", dbAddr, group, keyword, i*50)
				resp, err := httpGet(searchAddr)
				if err != nil {
					return err
				}
				if err := s.parse(resp); err != nil {
					return err
				}
				time.Sleep(sleepTime * time.Second)
			}
		}
	}
	return nil
}

func callback(i int, s *goquery.Selection) {
	title, _ := s.Find("a").Attr("title")
	href, _ := s.Find("a").Attr("href")
	time, _ := s.Find(".td-time").Attr("title")
	fmt.Printf("Review %d: %s %s %s\n", i, title, href, time)
	for _, filterKeyword := range filterKeywords {
		if strings.Contains(title, filterKeyword) {
			log.Println("过滤关键词:", filterKeyword)
			return
		}
	}
	data := fmt.Sprintf("<div><a href=%s>%s</a> %s</div>", href, title, time)
	result += data
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
	err := ioutil.WriteFile("result.html", []byte(result), 0644)
	if err != nil {
		log.Println(err)
	}
	cmdstr := "open result.html"
	cmd := exec.Command("bash", "-c", cmdstr)
	_, err = cmd.CombinedOutput()
	if err != nil {
		log.Println("cmd:", cmdstr, "err:", err)
	}
}
