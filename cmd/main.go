package main

import (
	"fmt"
	"log"

	"github.com/my-Sakura/go-spider-scheduler/model/csv"
)

// func main() {
// 	url := "https://piao.qunar.com/ticket/list.htm?keyword=%E4%BF%9D%E5%AE%9A&region=%E6%B2%B3%E5%8C%97&from=mps_search_suggest&page=1"
// 	st := scenicTicket.New()
// 	err := st.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}
// 	err = csv.WriteIntoCSVFile("scenicTicket.csv", st.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	url := "http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A"
// 	bdsc := selectedContent.NewBaoDingSelectedContentSummary()
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
// 	bdsc := selectedContent.NewBeiJingSelectedContentSummary()
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
// 	bdpt := popularTourist.NewBeiJingPopularTouristSummary()
// 	err := bdpt.Crawl(url)
// 	if err != nil {
// 		log.Println(err)
// 	}

// 	err = csv.WriteIntoCSVFile("beiJingPopularTourist.csv", bdpt.Slice())
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func main() {
	url := "http://www.mafengwo.cn/search/q.php?q=%E4%BF%9D%E5%AE%9A"
	bdpt := popularTourist.NewBeiJingPopularTouristSummary()
	err := bdpt.Crawl(url)
	if err != nil {
		log.Println(err)
	}

	err = csv.WriteIntoCSVFile("tes.csv", bdpt.Slice())
	if err != nil {
		fmt.Println(err)
	}
}
