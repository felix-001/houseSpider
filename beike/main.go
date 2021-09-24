package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	ErrStatusCode = errors.New("status code err")
	ErrParseHtml  = errors.New("parse html err")
)

func httpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", ErrStatusCode
	}
	return string(body), nil
}

func getData(raw, key string) (string, error) {
	start := strings.Index(raw, key)
	if start == -1 {
		log.Println("parse html error")
		return "", ErrParseHtml
	}
	new := raw[start+len(key):]
	end := strings.Index(new, " ")
	if end == -1 {
		log.Println("parse html error")
		return "", ErrParseHtml
	}
	res := new[:end]
	return res, nil
}

func getHouseData() (string, string, error) {
	body, err := httpGet("https://bj.ke.com/")
	if err != nil {
		return "", "", err
	}
	secondhand, err := getData(body, "北京在售二手房 ")
	if err != nil {
		return "", "", err
	}
	new, err := getData(body, "北京在售新房楼盘 ")
	if err != nil {
		return "", "", err
	}
	return secondhand, new, nil
}

func appendDataToCSV(secondhand, new string) error {
	file := "./data.csv"
	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	fi, err := os.Stat(file)
	if err != nil {
		log.Println(err)
		return err
	}
	if fi.Size() == 0 {
		_, err = io.WriteString(f, "日期, 在售二手房, 在售新房楼盘\n")
		if err != nil {
			log.Println(err)
			return err
		}
	}
	defer f.Close()
	date := time.Now().Format("2006-01-02 15:04:05")
	s := date + ", " + secondhand + ", " + new + "\n"
	_, err = io.WriteString(f, s)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func main() {
	log.SetFlags(log.Lshortfile)
	secondhand, new, err := getHouseData()
	if err != nil {
		return
	}
	if err := appendDataToCSV(secondhand, new); err != nil {
		return
	}
	log.Println(secondhand, new)
}
