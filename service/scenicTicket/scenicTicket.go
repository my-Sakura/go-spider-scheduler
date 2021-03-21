package scenicTicket

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/my-Sakura/go-spider-scheduler/service/request"
)

type ScenicTicketSummary struct {
	Data        []ScenicTicketItem
	tableHeader []string
	client      *http.Client
}

type ScenicTicketItem struct {
	TicketName   string
	TicketPrice  int
	MonthlySales int
}

// New
func New() *ScenicTicketSummary {
	return &ScenicTicketSummary{
		Data:   make([]ScenicTicketItem, 0),
		client: &http.Client{},
	}
}

// Crawl -
func (s *ScenicTicketSummary) Crawl(URL string) error {
	req := request.Get(URL, nil)
	if req == nil {
		return request.ErrInvalidRequest
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}

	doc.Find("div.sight_item").Each(func(i int, selection *goquery.Selection) {
		scenic, exist := selection.Attr("data-sight-name")
		if !exist {
			log.Println("ticketName crawl fail")
		}
		ticketPrice := selection.Find("div > div.sight_item_pop > table > tbody > tr:nth-child(1) > td > span > em").Text()
		monthlySales := selection.Find("div > div.sight_item_pop > table > tbody > tr:nth-child(4) > td > span").Text()
		price, err := strconv.Atoi(ticketPrice)
		if err != nil {
			log.Printf("ticketPrice string convert to int failed")
		}
		monthly, err := strconv.Atoi(monthlySales)
		if err != nil {
			log.Printf("monthly string convert to int failed")
		}
		scenicTicketItem := ScenicTicketItem{
			TicketName:   scenic,
			TicketPrice:  price,
			MonthlySales: monthly,
		}

		s.Data = append(s.Data, scenicTicketItem)
	})

	return nil
}

// Slice -
func (s *ScenicTicketSummary) Slice() [][]string {
	scenicTicketSlices := make([][]string, 0)
	scenicTicketSlice := make([]string, 0)

	s.setTableHeader()
	scenicTicketSlices = append(scenicTicketSlices, s.tableHeader)
	for _, data := range s.Data {
		scenicTicketSlice = append(scenicTicketSlice, data.TicketName)
		scenicTicketSlice = append(scenicTicketSlice, strconv.Itoa(data.TicketPrice))
		scenicTicketSlice = append(scenicTicketSlice, strconv.Itoa(data.MonthlySales))
		scenicTicketSlices = append(scenicTicketSlices, scenicTicketSlice)
		scenicTicketSlice = nil
	}

	return scenicTicketSlices
}

func (s *ScenicTicketSummary) setTableHeader() {
	s.tableHeader = append(s.tableHeader, "ticketName", "tickerPrice", "monthlySales")
}
