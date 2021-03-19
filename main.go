package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	popularTourists  = make(map[string][]PopularTourist)
	selectedContents = make(map[string][]SelectedContent)
)

type PopularTourist struct {
	TouristName string
	PublishDate string
	CrawlDate   string
	HotComments []HotComment
}

type HotComment struct {
	CommentatorName string
	CommentContent  string
}

type SelectedContent struct {
	Title         string
	PublishedDate string
	CrawlDate     string
	ScanNumber    int
}

type Task struct {
	client *http.Client
	mu     *sync.Mutex
}

func NewTask() *Task {
	return &Task{
		client: &http.Client{},
		mu:     &sync.Mutex{},
	}
}

func (t *Task) crawlManChengHanMu(URL string) {
	// var selectedContent SelectedContent
	// selectedContentSlice := make([]SelectedContent, 0)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// doc.Find()
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	datas := make([][]string, 0)
	data := make([]string, 0)
	task := NewTask()

	task.crawlManChengHanMu("https://piao.qunar.com/ticket/detail_3376839007.html?st=a3clM0QlRTQlQkYlOUQlRTUlQUUlOUElMjZpZCUzRDM1MzklMjZ0eXBlJTNEMCUyNmlkeCUzRDElMjZxdCUzRHJlZ2lvbiUyNmFwayUzRDIlMjZzYyUzRFdXVyUyNnVyJTNEJUU2JUIyJUIzJUU1JThDJTk3JTI2bHIlM0QlRTQlQkYlOUQlRTUlQUUlOUElMjZmdCUzRCU3QiU3RA%3D%3D#from=mps_search_suggest")

	go func() {
		task.crawlBeiJingSelectedContent("http://www.mafengwo.cn/search/q.php?q=%E5%8C%97%E4%BA%AC")
		for _, beiJingSelectedContent := range selectedContents["beijing"] {
			data = append(data, beiJingSelectedContent.Title)
			data = append(data, strconv.Itoa(beiJingSelectedContent.ScanNumber))
			data = append(data, beiJingSelectedContent.PublishedDate)
			data = append(data, beiJingSelectedContent.CrawlDate)
			datas = append(datas, data)
			data = nil
		}
		writeIntoCSVFile("beiJingSelected.csv", datas)
		datas = nil
		wg.Done()
	}()

	go func() {
		task.crawlBaoDingSelectedContent("http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A")
		for _, baoDingSelectedContent := range selectedContents["baoding"] {
			data = append(data, baoDingSelectedContent.Title)
			data = append(data, strconv.Itoa(baoDingSelectedContent.ScanNumber))
			data = append(data, baoDingSelectedContent.PublishedDate)
			data = append(data, baoDingSelectedContent.CrawlDate)
			datas = append(datas, data)
			data = nil
		}
		writeIntoCSVFile("baoDingSelected.csv", datas)
		datas = nil
		wg.Done()
	}()

	go func() {
		task.crawlBaoDingPopularTourist("http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A")
		for _, baoDingPopularTourist := range popularTourists["baoding"] {
			data = append(data, baoDingPopularTourist.TouristName)
			data = append(data, baoDingPopularTourist.PublishDate)
			data = append(data, baoDingPopularTourist.CrawlDate)
			datas = append(datas, data)
			data = nil
		}
		writeIntoCSVFile("baoDingPopularTourist.csv", datas)
		datas = nil
		wg.Done()
	}()

	go func() {
		task.crawlBeiJingPopularTourist("http://www.mafengwo.cn/search/q.php?q=%E5%8C%97%E4%BA%AC")
		for _, beiJingPopularTourist := range popularTourists["beijing"] {
			data = append(data, beiJingPopularTourist.TouristName)
			data = append(data, beiJingPopularTourist.PublishDate)
			data = append(data, beiJingPopularTourist.CrawlDate)
			datas = append(datas, data)
			data = nil
		}
		writeIntoCSVFile("beiJingPopularTourist.csv", datas)
		datas = nil
		wg.Done()
	}()

	wg.Wait()
}

func (t *Task) crawlBaoDingSelectedContent(URL string) {
	var selectedContent SelectedContent
	selectedContentSlice := make([]SelectedContent, 0)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	regexpScanNumber := regexp.MustCompile("[0-9]+浏览")
	links := make([]string, 0)

	doc.Find("#_j_search_result_left > div:nth-child(5) > ul > li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.foot").Text(), "浏览") {
			selectedContent.Title = s.Find("span.head").Text()
			selectedContent.CrawlDate = time.Now().Format("2006-01-02")
			selectedContent.PublishedDate = regexpDate.FindString(s.Find("div.foot").Text())
			selectedContent.ScanNumber, _ = strconv.Atoi(strings.Trim(regexpScanNumber.FindString(s.Find("div.foot").Text()), "浏览"))
			selectedContentSlice = append(selectedContentSlice, selectedContent)
		} else {
			link, _ := s.Find("a[href]").Attr("href")
			links = append(links, link)
		}
	})

	for _, link := range links {
		req, err = http.NewRequest("GET", link, nil)
		resp, err = t.client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		t := doc.Find("#_js_askDetail > div.q-content > div.q-info1.clearfix > div.pub-bar.fr > span > span").Text()
		selectedContent.Title = doc.Find("#_js_askDetail > div.q-content > div.q-title > h1 > a").Text()
		selectedContent.CrawlDate = time.Now().Format("2006-01-02")
		selectedContent.PublishedDate = regexpDate.FindString(t)
		selectedContent.ScanNumber, _ = strconv.Atoi(strings.Trim(doc.Find("#_js_askDetail > div.q-operate.clearfix > div.fr > span:nth-child(1)").Text(), "浏览"))
		selectedContentSlice = append(selectedContentSlice, selectedContent)
	}

	t.mu.Lock()
	if _, ok := selectedContents["baoding"]; ok {
		delete(selectedContents, "baoding")
	}

	selectedContents["baoding"] = selectedContentSlice
	t.mu.Unlock()
}

func (t *Task) crawlBeiJingSelectedContent(URL string) {
	var selectedContent SelectedContent
	selectedContentSlice := make([]SelectedContent, 0)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")
	// req.Header.Set("Cookie", cookie)

	resp, err := t.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	regexpScanNumber := regexp.MustCompile("[0-9]+浏览")
	links := make([]string, 0)

	doc.Find("#_j_search_result_left > div:nth-child(7) > ul > li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.foot").Text(), "浏览") {
			selectedContent.Title = s.Find("span.head").Text()
			selectedContent.CrawlDate = time.Now().Format("2006-01-02")
			selectedContent.PublishedDate = regexpDate.FindString(s.Find("div.foot").Text())
			selectedContent.ScanNumber, _ = strconv.Atoi(strings.Trim(regexpScanNumber.FindString(s.Find("div.foot").Text()), "浏览"))
			selectedContentSlice = append(selectedContentSlice, selectedContent)
		} else {
			link, _ := s.Find("a[href]").Attr("href")
			links = append(links, link)
		}
	})

	for _, link := range links {
		req, err = http.NewRequest("GET", link, nil)
		resp, err = t.client.Do(req)
		if err != nil {
			log.Println(err)
		}

		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		t := doc.Find("#_js_askDetail > div.q-content > div.q-info1.clearfix > div.pub-bar.fr > span > span").Text()
		selectedContent.Title = doc.Find("#_js_askDetail > div.q-content > div.q-title > h1 > a").Text()
		selectedContent.CrawlDate = time.Now().Format("2006-01-02")
		selectedContent.PublishedDate = regexpDate.FindString(t)
		selectedContent.ScanNumber, _ = strconv.Atoi(strings.Trim(doc.Find("#_js_askDetail > div.q-operate.clearfix > div.fr > span:nth-child(1)").Text(), "浏览"))
		selectedContentSlice = append(selectedContentSlice, selectedContent)
	}

	t.mu.Lock()
	if _, ok := selectedContents["beijing"]; ok {
		delete(selectedContents, "beijing")
	}

	selectedContents["beijing"] = selectedContentSlice
	t.mu.Unlock()
}

func (t *Task) crawlBaoDingPopularTourist(URL string) {
	var popularTourist PopularTourist
	popularTouristSlice := make([]PopularTourist, 0)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	links := make([]string, 0)
	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	regexpLink := regexp.MustCompile("&id=[0-9]+")

	doc.Find("#_j_search_result_left > div:nth-child(3) > div.content.top_pois-list > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id := strings.Trim(regexpLink.FindString(href), "&id=")
		host := "http://mafengwo.cn"
		path := "poi" + "/" + id + ".html"
		link := host + "/" + path
		links = append(links, link)
	})

	for _, link := range links {
		req.URL, _ = url.Parse(link)
		resp, err := t.client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}

		popularTourist.TouristName = doc.Find("body > div.container > div.row.row-top > div > div.title > h1").Text()
		popularTourist.PublishDate = regexpDate.FindString(doc.Find("body > div.container > div:nth-child(7) > div.mod.mod-detail > div:nth-child(6)").Text())
		popularTourist.CrawlDate = time.Now().Format("2006-01-02")
		doc.Find("#pagelet-block-24289b4cb9321822e98f07ac1c91450d > div > div._j_commentlist > div.rev-list > ul > li").Each(func(i int, s *goquery.Selection) {
			comment := s.Find("p").Text()
			fmt.Println(comment)
		})
		popularTouristSlice = append(popularTouristSlice, popularTourist)
	}

	t.mu.Lock()
	if _, ok := popularTourists["baoding"]; ok {
		delete(popularTourists, "baoding")
	}
	popularTourists["baoding"] = popularTouristSlice
	t.mu.Unlock()
}

func (t *Task) crawlBeiJingPopularTourist(URL string) {
	var popularTourist PopularTourist
	popularTouristSlice := make([]PopularTourist, 0)

	req, _ := http.NewRequest("GET", URL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := t.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	links := make([]string, 0)
	regexpLink := regexp.MustCompile("&id=[0-9]+")
	doc.Find("#_j_search_result_left > div:nth-child(3) > div.content.top_pois-list > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id := strings.Trim(regexpLink.FindString(href), "&id=")
		host := "http://mafengwo.cn"
		path := "poi" + "/" + id + ".html"
		link := host + "/" + path
		links = append(links, link)
	})

	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	for _, link := range links {
		req.URL, _ = url.Parse(link)
		resp, err = t.client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}
		popularTourist.TouristName = doc.Find("body > div.container > div.row.row-top > div > div.title > h1").Text()
		popularTourist.PublishDate = regexpDate.FindString(doc.Find("body > div.container > div:nth-child(7) > div.mod.mod-detail > div:nth-child(6)").Text())
		popularTourist.CrawlDate = time.Now().Format("2006-01-02")
		popularTouristSlice = append(popularTouristSlice, popularTourist)
	}

	t.mu.Lock()
	if _, ok := popularTourists["beijing"]; ok {
		delete(popularTourists, "beijing")
	}
	popularTourists["beijing"] = popularTouristSlice
	t.mu.Unlock()
}

func writeIntoCSVFile(fileName string, datas [][]string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}

	w := csv.NewWriter(file)
	defer w.Flush()
	for _, data := range datas {
		w.Write(data)
	}
}
