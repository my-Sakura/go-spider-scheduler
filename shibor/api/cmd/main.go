package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/my-Sakura/go-spider-scheduler/api/service/structure/shibor"
)

func main() {
	fmt.Println("start")
	http.HandleFunc("/api/v1/shiborbid", ShiborBid)
	http.HandleFunc("/api/v1/shibor", Shibor)
	http.HandleFunc("/api/v1/shiboraverage", ShiborAverage)
	http.ListenAndServe(":1002", nil)
}

func ShiborBid(w http.ResponseWriter, r *http.Request) {
	url := "http://www.chinamoney.com.cn/ags/ms/cm-u-bk-shibor/ShiborPri"
	param := r.URL.Query()
	url += "?" + param.Encode()

	shiborBidSummary := shibor.NewShiborBidSummary()

	data, err := shiborBidSummary.Crawl(url)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data faied"))
	}

	err = shiborBidSummary.Parse(data)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data failed"))
	}

	resp, err := json.Marshal(shiborBidSummary.Records)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	w.Write(resp)
}

func Shibor(w http.ResponseWriter, r *http.Request) {
	url := "http://www.chinamoney.com.cn/r/cms/www/chinamoney/data/shibor/shibor.json"

	shiborSummary := shibor.NewShiborSummary()

	data, err := shiborSummary.Crawl(url)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data faied"))
	}

	err = shiborSummary.Parse(data)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data failed"))
	}

	resp, err := json.Marshal(shiborSummary.Records)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	w.Write(resp)
}

func ShiborAverage(w http.ResponseWriter, r *http.Request) {
	url := "http://www.chinamoney.com.cn/r/cms/www/chinamoney/data/shibor/shibor-mn.json"

	shiborAverageSummary := shibor.NewShiborAverageSummary()

	data, err := shiborAverageSummary.Crawl(url)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data faied"))
	}

	err = shiborAverageSummary.Parse(data)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("get data failed"))
	}

	resp, err := json.Marshal(shiborAverageSummary.Records)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
	}

	w.Write(resp)
}
