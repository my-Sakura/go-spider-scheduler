package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	nats "github.com/nats-io/nats.go"
)

func init() {
	domainCollector["http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/"] = initAcademicianCollector()
	domainCollector["http://www.gov.cn/xinwen/yaowen.htm"] = initStateDepartmentNewsCrawl()
	domainCollector["http://www.gov.cn/zhengce/index.htm"] = initStateDepartmentPoliciesCrawl()
	var err error
	nc, err = nats.Connect(natsURL)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	urls := []string{"http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/", "http://www.gov.cn/xinwen/yaowen.htm", "http://www.gov.cn/zhengce/index.htm"}
	for _, url := range urls {
		instance := factory(url)
		instance.Visit(url)
	}
}

var domainCollector = map[string]*colly.Collector{}
var nc *nats.Conn
var natsURL = "nats://localhost:4222"

func factory(urlStr string) *colly.Collector {
	return domainCollector[urlStr]
}

func initAcademicianCollector() *colly.Collector {
	reg := regexp.MustCompile("[1-9]+")

	var count int
	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
	)

	c.OnHTML("dl#allNameBar", func(e *colly.HTMLElement) {
		e.ForEach("dt", func(_ int, element *colly.HTMLElement) {
			department_people_number, _ := strconv.Atoi(reg.FindString(element.ChildText("em")))

			for i := count; i < department_people_number+count; i++ {
				link, _ := e.DOM.Find("a[href]").Eq(i).Attr("href")

				c.Visit(element.Request.AbsoluteURL(link))
			}
			count += department_people_number
		})
		nc.Flush()
	})

	c.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.111:7890"
		if r.URL.String() != "http://www.casad.cas.cn/ysxx2017/ysmdyjj/qtysmd_124280/" {
			err := nc.Publish("academician", []byte(r.URL.String()))
			if err != nil {
				log.Println("academicianURL publish fail")
			}
		}
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	return c
}

func initStateDepartmentNewsCrawl() *colly.Collector {
	now := time.Now()
	date := now.Format("2006-01-02")

	c := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36"),
	)

	c.OnHTML("div.news_box", func(e *colly.HTMLElement) {
		e.ForEach("h4", func(_ int, element *colly.HTMLElement) {
			if element.ChildText("span.date") == date {
				link := element.ChildAttr("a[href]", "href")

				c.Visit(element.Request.AbsoluteURL(link))
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		r.ProxyURL = "http://192.168.0.111:7890"
		if r.URL.String() != "http://www.gov.cn/xinwen/yaowen.htm" {
			err := nc.Publish("stateDepartmentNews", []byte(r.URL.String()))
			if err != nil {
				log.Println("stateDepartmentNewsURL publish fail")
			}
		}
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	return c
}

func initStateDepartmentPoliciesCrawl() *colly.Collector {
	now := time.Now()
	date := now.Format("2006-01-02")
	c := colly.NewCollector()

	c.OnHTML("div.list_left_con", func(e *colly.HTMLElement) {
		e.ForEach("div.latestPolicy_left_item", func(_ int, element *colly.HTMLElement) {
			if element.DOM.Find("span").Eq(0).Text() == date {
				link := element.ChildAttr("a[href]", "href")

				c.Visit(e.Request.AbsoluteURL(link))
			}
		})
	})

	c.OnRequest(func(r *colly.Request) {
		if r.URL.String() != "http://www.gov.cn/zhengce/index.htm" {
			err := nc.Publish("stateDepartmentPoliciesCrawl", []byte(r.URL.String()))
			if err != nil {
				log.Println("stateDepartmentPoliciesCrawlURL publish fail")
			}
		}
		fmt.Printf("visiting => %s\n", r.URL.String())
	})

	return c
}
