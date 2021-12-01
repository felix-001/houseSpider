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
	"os/exec"
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
	Title   = "北京全区域二手房源走势图"
	CsvFile = "/Users/rigensen/workspace/learn/houseSpider/beike/data.csv"
	PngPath = "/Users/rigensen/workspace/learn/houseSpider/beike"
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

func getRegionData(raw string) (string, error) {
	start := strings.Index(raw, "共找到")
	if start == -1 {
		log.Println("parse html error")
		return "", ErrParseHtml
	}
	start += len("共找到")
	new := raw[start:]
	start = strings.Index(new, "<span> ")
	if start == -1 {
		log.Println("parse html error")
		return "", ErrParseHtml
	}
	start += len("<span> ")
	new = new[start:]
	end := strings.Index(new, " </span>")
	if end == -1 {
		log.Println("parse html error")
		return "", ErrParseHtml
	}
	log.Println(new[:end])
	return new[:end], nil
}

func getRegionHouseData(regions []string) ([]string, error) {
	nums := []string{}
	for _, region := range regions {
		addr := fmt.Sprintf("https://bj.ke.com/ershoufang/rs%s/", region)
		body, err := httpGet(addr)
		if err != nil {
			return nil, err
		}
		num, err := getRegionData(body)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func appendDataToCSV(secondhand, new string, regionNums []string) error {
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
	s := date + ", " + secondhand + ", " + new
	for _, num := range regionNums {
		s += ", " + num
	}
	s += "\n"
	_, err = io.WriteString(f, s)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func savePNG(title, xLabel, yLabel string, xys plotter.XYs) error {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = xLabel
	p.Y.Label.Text = yLabel

	err := plotutil.AddLinePoints(p, Second, xys)
	if err != nil {
		log.Println(err)
		return err
	}

	pngFile := fmt.Sprintf("%s/%d.png", PngPath, title)
	if err := p.Save(8*vg.Inch, 8*vg.Inch, pngFile); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func parseCSV() ([]plotter.XYs, error) {
	file, err := os.Open(CsvFile)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	xys := make([]plotter.XYs, 10)
	for idx := range xys {
		xys[idx] = make(plotter.XYs, 0)
	}
	i := 0
	for scanner.Scan() {
		if i == 0 {
			i++
			continue
		}
		items := strings.Split(scanner.Text(), ",")
		idx := 0
		//log.Println(len(items))
		for _, item := range items {
			if idx == 0 {
				idx++
				continue
			}
			y, err := strconv.Atoi(strings.TrimSpace(item))
			if err != nil {
				log.Println(err)
				return nil, err
			}
			xy := plotter.XY{X: float64(i), Y: float64(y)}
			xys[idx] = append(xys[idx], xy)
			idx++
		}
		i++
	}
	return xys, nil
}

/*
func showPNG() {
	cmd := exec.Command("open", PngFile)
	cmd.Run()
}
*/

func uploadPNG() {
	cmd := exec.Command("bash", "-c", `cd /Users/rigensen/workspace/learn/houseSpider/beike; 
			     /usr/bin/git add .;
			     /usr/bin/git commit -m "update";
			     /usr/bin/git push`)
	b, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(string(b))

}

func main() {
	log.SetFlags(log.Lshortfile)
	secondhand, new, err := getHouseData()
	if err != nil {
		return
	}
	regions := []string{"东坝"}
	nums, err := getRegionHouseData(regions)
	if err != nil {
		return
	}
	if err := appendDataToCSV(secondhand, new, nums); err != nil {
		return
	}
	log.Println(secondhand, new)
	xys, err := parseCSV()
	if err != nil {
		return
	}
	log.Println(len(xys))
	if err := savePNG(Title, "time", "house count", xys[0]); err != nil {
		log.Println(err)
		return
	}
	//showPNG()
	uploadPNG()
}
