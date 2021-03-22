package service

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/my-Sakura/go-spider-scheduler/pkg/request"
)

type BeiJingSelectedContentSummary struct {
	Data        []BeiJingSelectedContentItem
	tableHeader []string
	client      *http.Client
}

type BeiJingSelectedContentItem struct {
	Title         string
	PublishedDate string
	CrawlDate     string
	ScanNumber    int
}

// NewBeiJingSelectedContentSummary -
func NewBeiJingSelectedContentSummary() *BeiJingSelectedContentSummary {
	return &BeiJingSelectedContentSummary{
		Data:   make([]BeiJingSelectedContentItem, 0),
		client: &http.Client{},
	}
}

func (b *BeiJingSelectedContentSummary) Crawl(URL string) error {
	req := request.Get(URL, nil)

	if req == nil {
		return request.ErrInvalidRequest
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := b.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	regexpScanNumber := regexp.MustCompile("[0-9]+浏览")
	links := make([]string, 0)
	doc.Find("#_j_search_result_left > div:nth-child(7) > ul > li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Find("div.foot").Text(), "浏览") {
			title := s.Find("span.head").Text()
			crawlDate := time.Now().Format("2006-01-02")
			publishedDate := regexpDate.FindString(s.Find("div.foot").Text())
			scanNumber, err := strconv.Atoi(strings.Trim(regexpScanNumber.FindString(s.Find("div.foot").Text()), "浏览"))
			if err != nil {
				log.Println("ScanNumber string convert to int failed")
			}
			beiJingSelectedContentItem := BeiJingSelectedContentItem{
				Title:         title,
				CrawlDate:     crawlDate,
				PublishedDate: publishedDate,
				ScanNumber:    scanNumber,
			}

			b.Data = append(b.Data, beiJingSelectedContentItem)
		} else {
			link, err := s.Find("a[href]").Attr("href")
			if !err {
				log.Println("link find error")
			}
			links = append(links, link)
		}
	})

	for _, link := range links {
		req := request.Get(link, nil)
		if req == nil {
			return request.ErrInvalidRequest
		}

		resp, err = b.client.Do(req)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return err
		}

		t := doc.Find("#_js_askDetail > div.q-content > div.q-info1.clearfix > div.pub-bar.fr > span > span").Text()
		title := doc.Find("#_js_askDetail > div.q-content > div.q-title > h1 > a").Text()
		crawlDate := time.Now().Format("2006-01-02")
		publishedDate := regexpDate.FindString(t)
		scanNumber, err := strconv.Atoi(strings.Trim(doc.Find("#_js_askDetail > div.q-operate.clearfix > div.fr > span:nth-child(1)").Text(), "浏览"))
		if err != nil {
			log.Println("ScanNumber string convert to int failed")
		}
		beiJingSelectedContentItem := BeiJingSelectedContentItem{
			Title:         title,
			CrawlDate:     crawlDate,
			PublishedDate: publishedDate,
			ScanNumber:    scanNumber,
		}
		b.Data = append(b.Data, beiJingSelectedContentItem)
	}

	return nil
}

// Slice -
func (s *BeiJingSelectedContentSummary) Slice() [][]string {
	BeiJingSelectedContentSlices := make([][]string, 0)
	BeiJingSelectedContentSlice := make([]string, 0)

	s.setTableHeader()
	BeiJingSelectedContentSlices = append(BeiJingSelectedContentSlices, s.tableHeader)
	for _, data := range s.Data {
		BeiJingSelectedContentSlice = append(BeiJingSelectedContentSlice, data.Title)
		BeiJingSelectedContentSlice = append(BeiJingSelectedContentSlice, data.PublishedDate)
		BeiJingSelectedContentSlice = append(BeiJingSelectedContentSlice, data.CrawlDate)
		BeiJingSelectedContentSlice = append(BeiJingSelectedContentSlice, strconv.Itoa(data.ScanNumber))
		BeiJingSelectedContentSlices = append(BeiJingSelectedContentSlices, BeiJingSelectedContentSlice)

		BeiJingSelectedContentSlice = nil
	}

	return BeiJingSelectedContentSlices
}

func (s *BeiJingSelectedContentSummary) setTableHeader() {
	s.tableHeader = append(s.tableHeader, "Title", "PublishedDate", "CrawlDate", "ScanNumber")
}
