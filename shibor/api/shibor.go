package api

import (
	"encoding/json"
	"net/http"

	"github.com/my-Sakura/go-spider-scheduler/shibor/api/service/structure/shibor"
)

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
