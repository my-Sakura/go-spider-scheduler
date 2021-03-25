package main

import (
	"fmt"
	"log"

	"github.com/my-Sakura/go-spider-scheduler/pkg/model/csv"
	"github.com/my-Sakura/go-spider-scheduler/pkg/service"
)

func main() {
	url := "https://piao.qunar.com/ticket/list.htm?keyword=%E4%BF%9D%E5%AE%9A&region=%E6%B2%B3%E5%8C%97&from=mps_search_suggest&page=1"
	st := service.New()
	err := st.Crawl(url)
	if err != nil {
		log.Println(err)
	}
	err = csv.WriteIntoCSVFile("scenicTicket.csv", st.Slice())
	if err != nil {
		fmt.Println(err)
	}
}

// func main() {
// 	url := "http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A"
// 	bdsc := service.NewBaoDingSelectedContentSummary()
// 	err := bdsc.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = csv.WriteIntoCSVFile("baoDingSelectedContent.csv", bdsc.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	url := "http://www.mafengwo.cn/search/q.php?q=%E5%8C%97%E4%BA%AC"
// 	bdsc := service.NewBeiJingSelectedContentSummary()
// 	err := bdsc.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = csv.WriteIntoCSVFile("beiJingSelectedContent.csv", bdsc.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	url := "http://www.mafengwo.cn/search/q.php?q=%E5%8C%97%E4%BA%AC"
// 	bdpt := service.NewBeiJingPopularTouristSummary()
// 	err := bdpt.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = csv.WriteIntoCSVFile("beiJingPopularTourist.csv", bdpt.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	url := "http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A&seid=72043C00-6DBB-485D-B694-8FD4EE536B9D"
// 	bdpt := service.NewBaoDingPopularTouristSummary()
// 	err := bdpt.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = csv.WriteIntoCSVFile("baoDingPopularTourist.csv", bdpt.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }
