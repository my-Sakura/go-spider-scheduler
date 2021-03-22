package api

import (
	"encoding/json"
	"net/http"

	"github.com/my-Sakura/go-spider-scheduler/shibor/api/service/structure/shibor"
)

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
