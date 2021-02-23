package main

import (
	"fmt"

	"github.com/gocolly/colly"
	nats "github.com/nats-io/nats.go"
)

func main() {
	startConsumer()
}

var natsURL = "nats://localhost:4222"

func academicianParser(url string) {
	c := colly.NewCollector()

	c.OnHTML("div.contentBar", func(e *colly.HTMLElement) {
		name := e.DOM.Find("h1").Eq(0).Text()
		content := e.DOM.Find("p:contains(院士)").Text()
		fmt.Println(name, content)
	})

	c.Visit(url)
}

func stateDepartmentParser(url string) {
	var content string
	var title string
	c := colly.NewCollector()
	fmt.Println(url, "in")
	c.OnHTML("div.content", func(e *colly.HTMLElement) {
		title = e.DOM.Find("h1").Text()
		fmt.Println(title)

		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		//avoid blank title
		if title != "" {
			fmt.Println(title, content, "have title")
		} else {
			fmt.Println(content, "have no title")
		}

		//clear content
		content = ""
	})

	c.Visit(url)
}

func stateDepartmentPoliciesParser(url string) {
	var title string
	var content string
	c := colly.NewCollector()

	c.OnHTML("div.article.oneColumn.pub_border", func(e *colly.HTMLElement) {
		title = e.DOM.Find("h1").Eq(0).Text()

		e.ForEach("div#UCAP-CONTENT.pages_content", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		fmt.Println(title, content)

		//clear content
		content = ""
	})

	c.OnHTML("td.b12c#UCAP-CONTENT", func(e *colly.HTMLElement) {
		e.ForEach("strong", func(_ int, element *colly.HTMLElement) {
			title = title + element.ChildText("span") + "\n"
		})

		e.ForEach("p", func(_ int, element *colly.HTMLElement) {
			content = content + element.Text + "\n"
		})

		fmt.Println(title, content)

		//clear content
		content = ""
	})

	c.Visit(url)
}

func startConsumer() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return
	}

	nc.Subscribe("academician", func(msg *nats.Msg) {
		urlStr := string(msg.Data)
		academicianParser(urlStr)
	})

	nc.Subscribe("stateDepartmentNews", func(msg *nats.Msg) {
		urlStr := string(msg.Data)
		stateDepartmentParser(urlStr)
	})

	nc.Subscribe("stateDepartmentPoliciesCrawl", func(msg *nats.Msg) {
		urlStr := string(msg.Data)
		stateDepartmentPoliciesParser(urlStr)
	})

	select {}
}
