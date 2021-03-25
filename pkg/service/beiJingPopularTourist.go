package service

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
	"github.com/my-Sakura/go-spider-scheduler/pkg/request"
)

type BeiJingPopularTouristSummary struct {
	Data        []BeiJingPopularTouristItem
	tableHeader []string
	client      *http.Client
}

type BeiJingPopularTouristItem struct {
	TouristName string
	PublishDate string
	CrawlDate   string
	HotComments []hotComment
}

// NewBaoDingPopularTouristSummary -
func NewBeiJingPopularTouristSummary() *BeiJingPopularTouristSummary {
	return &BeiJingPopularTouristSummary{
		Data:   make([]BeiJingPopularTouristItem, 0),
		client: &http.Client{},
	}
}

// Crawl -
func (b *BeiJingPopularTouristSummary) Crawl(URL string) error {
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
		host := "https://www.mafengwo.cn"
		path := "poi" + "/" + id + ".html"
		link := host + "/" + path
		poiID = append(poiID, id)
		links = append(links, link)
	})

	// cookie := os.Getenv("cookie")
	for i, link := range links {
		req.URL, err = url.Parse(link)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Referer", fmt.Sprintf("https://www.mafengwo.cn/poi/%s.html", poiID[i]))
		req.Header.Set("Cookie", "mfw_uuid=6051d21d-8f63-b7d3-9ead-41ebf15f1fab; __jsluid_h=25a7972432177208e8df8d486ed7a1be; __omc_chl=; uva=s%3A78%3A%22a%3A3%3A%7Bs%3A2%3A%22lt%22%3Bi%3A1615974943%3Bs%3A10%3A%22last_refer%22%3Bs%3A6%3A%22direct%22%3Bs%3A5%3A%22rhost%22%3Bs%3A0%3A%22%22%3B%7D%22%3B; __mfwurd=a%3A3%3A%7Bs%3A6%3A%22f_time%22%3Bi%3A1615974943%3Bs%3A9%3A%22f_rdomain%22%3Bs%3A0%3A%22%22%3Bs%3A6%3A%22f_host%22%3Bs%3A3%3A%22www%22%3B%7D; __mfwuuid=6051d21d-8f63-b7d3-9ead-41ebf15f1fab; __jsluid_s=9fc1078a0250b7a1b39a8bedfd7cae92; c=JQIejQoA-1616035651303-a003175e46e79-1010466491; _r=google; _rp=a%3A2%3A%7Bs%3A1%3A%22p%22%3Bs%3A15%3A%22www.google.com%2F%22%3Bs%3A1%3A%22t%22%3Bi%3A1616202699%3B%7D; __mfwothchid=referrer%7Cwww.google.com; __mfwc=referrer%7Cwww.google.com; oad_n=a%3A3%3A%7Bs%3A3%3A%22oid%22%3Bi%3A1029%3Bs%3A2%3A%22dm%22%3Bs%3A16%3A%22open.mafengwo.cn%22%3Bs%3A2%3A%22ft%22%3Bs%3A19%3A%222021-03-20+09%3A49%3A34%22%3B%7D; __omc_r=; Hm_lvt_8288b2ed37e5bc9b4c9f7008798d2de0=1616330424; UM_distinctid=17854cd03331da-05379e919366c1-6418207d-13c680-17854cd0338975; CNZZDATA30065558=cnzz_eid%3D1786931762-1616326218-%26ntime%3D1616419044; mfw_passport_redirect=https%3A%2F%2Fwww.mafengwo.cn%2Fpoi%2F3474.html; PHPSESSID=t8hb0i55bfhur1bnmb7u8vq5f7; _fmdata=Q3shBgO19TziiXI8st75aXyWG3%2BLLD2%2BKE1THqC7%2FeFUZfiQlj3jwBB61r3cwov33szT9iv5U7YfxTNz38DUcojqh2iEPor6w9nehkBtk5o%3D; _xid=SWniKVcpL7giHF9gUEJ8I8h0f1qQ1JmSbREyugZcJd0lfg9OKbkDR0N78Rlm5viBV5ye1Nq2m8GLzMnH0bz%2Fuw%3D%3D; Hm_lpvt_8288b2ed37e5bc9b4c9f7008798d2de0=1616427494; bottom_ad_status=1; __mfwlv=1616460557; __mfwvn=23; __mfwb=645ed5dcd685.1.direct; __mfwa=1615974943543.71078.34.1616460557727.1616463954483; __mfwlt=1616463954; __jsl_clearance_s=1616464392.979|0|bfQaa%2FJ%2FGGgVzClEVGo8YEgGRt4%3D")

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

		beiJingPopularTouristItem := BeiJingPopularTouristItem{
			TouristName: touristName,
			PublishDate: publishDate,
			CrawlDate:   crawlDate,
			HotComments: hotComments,
		}

		b.Data = append(b.Data, beiJingPopularTouristItem)
	}

	return nil
}

func (b *BeiJingPopularTouristSummary) crawlSelectedComment(req *http.Request) ([]hotComment, error) {
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
func (b *BeiJingPopularTouristSummary) Slice() [][]string {
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

func (b *BeiJingPopularTouristSummary) setTableHeader() {
	b.tableHeader = append(b.tableHeader, "TouristName", "PublishDate", "CrawlDate", "HotComments")
}
