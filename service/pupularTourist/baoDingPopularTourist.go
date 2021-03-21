package popularTourist

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/my-Sakura/go-spider-scheduler/service/request"
)

type BaoDingPopularTouristSummary struct {
	Data        []BaoDingPopularTouristItem
	tableHeader []string
	client      *http.Client
}

type BaoDingPopularTouristItem struct {
	TouristName string
	PublishDate string
	CrawlDate   string
	HotComments []hotComment
}

// NewBaoDingPopularTouristSummary -
func NewBaoDingPopularTouristSummary() *BaoDingPopularTouristSummary {
	return &BaoDingPopularTouristSummary{
		Data:   make([]BaoDingPopularTouristItem, 0),
		client: &http.Client{},
	}
}

// Crawl -
func (b *BaoDingPopularTouristSummary) Crawl(URL string) error {
	req := request.Get(URL, nil)
	if req == nil {
		return request.ErrInvalidRequest
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := b.client.Do(req)
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
	poiID := make([]string, 0)
	regexpDate := regexp.MustCompile("[0-9]+-[0-9]+-[0-9]+")
	regexpLink := regexp.MustCompile("&id=[0-9]+")

	doc.Find("#_j_search_result_left > div:nth-child(3) > div.content.top_pois-list > a").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		id := strings.Trim(regexpLink.FindString(href), "&id=")
		host := "http://mafengwo.cn"
		path := "poi" + "/" + id + ".html"
		link := host + "/" + path
		poiID = append(poiID, id)
		links = append(links, link)
	})

	for i, link := range links {
		req.Header.Set("Referer", fmt.Sprintf("https://www.mafengwo.cn/poi/%s.html", poiID[i]))
		req.URL, _ = url.Parse(link)
		resp, err := b.client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			panic(err)
		}

		touristName := doc.Find("body > div.container > div.row.row-top > div > div.title > h1").Text()
		publishDate := regexpDate.FindString(doc.Find("body > div.container > div:nth-child(7) > div.mod.mod-detail > div:nth-child(6)").Text())
		crawlDate := time.Now().Format("2006-01-02")

		link = fmt.Sprintf("https://pagelet.mafengwo.cn/poi/pagelet/poiCommentListApi?params={\"poi_id\":\"%s\",\"page\":%d}", poiID[i], 1)
		r := request.Post(link, nil)
		if r == nil {
			return request.ErrInvalidRequest
		}

		r.Header.Set("Referer", fmt.Sprintf("https://www.mafengwo.cn/poi/%s.html", poiID[i]))
		r.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

		hotComments, err := b.crawlSelectedComment(r)
		if err != nil {
			return err
		}

		baoDingPopularTouristItem := BaoDingPopularTouristItem{
			TouristName: touristName,
			PublishDate: publishDate,
			CrawlDate:   crawlDate,
			HotComments: hotComments,
		}

		b.Data = append(b.Data, baoDingPopularTouristItem)
	}

	return nil
}

func (b *BaoDingPopularTouristSummary) crawlSelectedComment(req *http.Request) ([]hotComment, error) {
	comment := Comment{}
	hotComments := make([]hotComment, 0)

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, errors.New("baoDingHotComment crawl failed")
	}

	err = json.NewDecoder(resp.Body).Decode(&comment)
	if err != nil {
		return nil, errors.New("hotComment parse fail")
	}

	defer resp.Body.Close()

	reg := regexp.MustCompile("<p class=\"rev-txt\">[\\s\\S]*?</p>")
	result := reg.FindAllString(comment.Data.Html, -1)
	for _, v := range result {
		v = strings.Trim(v, "<p class=\"rev-txt\">")
		v = strings.Trim(v, "</p>")
		v = strings.ReplaceAll(v, "<br />", "\n")
		hotComment := hotComment{
			CommentContent: v,
		}
		hotComments = append(hotComments, hotComment)
	}

	return hotComments, nil
}

// Slice -
func (b *BaoDingPopularTouristSummary) Slice() [][]string {
	BaoDingSelectedContentSlices := make([][]string, 0)
	BaoDingSelectedContentSlice := make([]string, 0)

	b.setTableHeader()
	BaoDingSelectedContentSlices = append(BaoDingSelectedContentSlices, b.tableHeader)
	for _, data := range b.Data {
		BaoDingSelectedContentSlice = append(BaoDingSelectedContentSlice, data.TouristName)
		BaoDingSelectedContentSlice = append(BaoDingSelectedContentSlice, data.PublishDate)
		BaoDingSelectedContentSlice = append(BaoDingSelectedContentSlice, data.CrawlDate)
		BaoDingSelectedContentSlice = append(BaoDingSelectedContentSlice, sliceToString(data.HotComments))

		BaoDingSelectedContentSlices = append(BaoDingSelectedContentSlices, BaoDingSelectedContentSlice)

		BaoDingSelectedContentSlice = nil
	}

	return BaoDingSelectedContentSlices
}

func (b *BaoDingPopularTouristSummary) setTableHeader() {
	b.tableHeader = append(b.tableHeader, "TouristName", "PublishDate", "CrawlDate", "HotComments")
}
