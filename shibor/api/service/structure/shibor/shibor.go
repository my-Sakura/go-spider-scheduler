package shibor

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/my-Sakura/go-spider-scheduler/api/service"
	"github.com/my-Sakura/go-spider-scheduler/pkg/service/request"
)

type ShiborSummary struct {
	Records []ShiborDetail `json:"records"`
}

type ShiborDetail struct {
	TermCode        string  `json:"termCode"`
	Shibor          string  `json:"shibor"`
	ShibIdUpDown    string  `json:"shibIdUpDown"`
	ShibIdUpDownNum float64 `json:"shibIdUpDownNum"`
	TermCodePath    string  `json:"termCodePath"`
}

type ShiborAverageSummary struct {
	Records []ShiborAverageDetail `json:"records"`
}

type ShiborAverageDetail struct {
	TermCode   string   `json:"termCode"`
	ActionCode string   `json:"actionCode"`
	List       []string `json:"list"`
}

type ShiborBidSummary struct {
	Records []ShiborBidDetail `json:"records"`
	Data    ShiborBidData     `json:"data"`
}

type ShiborBidData struct {
	// Abbrname is name of the bank
	Abbrname string `json:"abbrname"`
}

type ShiborBidDetail struct {
	TermCode string `json:"termCode"`
	Bid      string `json:"bid"`
}

func NewShiborSummary() *ShiborSummary {
	return &ShiborSummary{
		Records: make([]ShiborDetail, 0),
	}
}

func (s *ShiborSummary) Crawl(URL string) (string, error) {
	req := request.Get(URL, nil)
	if req == nil {
		return "", service.ErrInvalidRequest
	}

	req.Header.Set("Referer", "http://www.chinamoney.com.cn/chinese/bkshibor/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (s *ShiborSummary) Parse(data string) error {
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		return err
	}

	return nil
}

func NewShiborAverageSummary() *ShiborAverageSummary {
	return &ShiborAverageSummary{
		Records: make([]ShiborAverageDetail, 0),
	}
}

func (s *ShiborAverageSummary) Crawl(URL string) (string, error) {
	req := request.Get(URL, nil)
	if req == nil {
		return "", service.ErrInvalidRequest
	}

	req.Header.Set("Referer", "http://www.chinamoney.com.cn/chinese/bkshibor/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (s *ShiborAverageSummary) Parse(data string) error {
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		return err
	}

	return nil
}

func NewShiborBidSummary() *ShiborBidSummary {
	return &ShiborBidSummary{
		Records: make([]ShiborBidDetail, 0),
	}
}

func (s *ShiborBidSummary) Crawl(URL string) (string, error) {
	req := request.Get(URL, nil)
	if req == nil {
		return "", service.ErrInvalidRequest
	}

	req.Header.Set("Referer", "http://www.chinamoney.com.cn/chinese/bkshibor/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.82 Safari/537.36")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (s *ShiborBidSummary) Parse(data string) error {
	err := json.Unmarshal([]byte(data), s)
	if err != nil {
		return err
	}

	return nil
}
