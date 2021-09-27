package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var (
	ErrStatusCode = errors.New("status code err")
	ErrParseHtml  = errors.New("parse html err")
)

const (
	Second  = "北京在售二手房"
	New     = "北京在售新房楼盘"
	Title   = "北京房源走势图"
	CsvFile = "/Users/rigensen/workspace/learn/houseSpider/beike/data.csv"
	PngFile = "/Users/rigensen/workspace/learn/houseSpider/beike/output.png"
	Url     = "https://bj.ke.com/"
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
	body, err := httpGet(Url)
	if err != nil {
		return "", "", err
	}
	secondhand, err := getData(body, Second+" ")
	if err != nil {
		return "", "", err
	}
	new, err := getData(body, New+" ")
	if err != nil {
		return "", "", err
	}
	return secondhand, new, nil
}

func appendDataToCSV(secondhand, new string) error {
	f, err := os.OpenFile(CsvFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	fi, err := os.Stat(CsvFile)
	if err != nil {
		log.Println(err)
		return err
	}
	if fi.Size() == 0 {
		_, err = io.WriteString(f, "日期, "+Second+", "+New+"\n")
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

func savePNG(secondhands, news plotter.XYs) error {
	p := plot.New()

	p.Title.Text = Title
	p.X.Label.Text = "time"
	p.Y.Label.Text = "house count"

	err := plotutil.AddLinePoints(p, Second, secondhands)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := p.Save(8*vg.Inch, 8*vg.Inch, PngFile); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func parseCSV() (plotter.XYs, plotter.XYs, error) {
	file, err := os.Open(CsvFile)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	seconds := plotter.XYs{}
	i := 0
	for scanner.Scan() {
		if i > 0 {
			items := strings.Split(scanner.Text(), ",")
			y, err := strconv.Atoi(strings.TrimSpace(items[1]))
			if err != nil {
				log.Println(err)
				return nil, nil, err
			}
			second := plotter.XY{X: float64(i), Y: float64(y)}
			seconds = append(seconds, second)
		}
		i++
	}
	return seconds, nil, nil
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
	seconds, news, err := parseCSV()
	if err != nil {
		return
	}
	log.Println(seconds, news)
	if err := savePNG(seconds, news); err != nil {
		log.Println(err)
		return
	}
}
