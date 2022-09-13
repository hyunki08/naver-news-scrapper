package scrapper

import (
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type newsItem struct {
	title     string
	desc      string
	publisher string
	pubDate   string
	link      string
}

func Scrapper(query string, maxPage int) {
	var newsList []newsItem

	baseUrl := "https://search.naver.com/search.naver"
	params := "?where=news&sm=tab_opt&sort=0&photo=3&field=0&pd=3&docid=&related=0&mynews=0&office_type=&office_section_code=&news_office_checked=&is_sug_officeid=0"
	params += ("&query=" + url.QueryEscape(query))
	now := time.Now()
	params += ("&de=" + now.Format("2006.01.02"))
	now = time.Now().AddDate(-5, 0, 0)
	params += ("&ds=" + now.Format("2006.01.02"))

	for i := 0; i < maxPage; i++ {
		l := getPage(baseUrl, params, i)
		newsList = append(l, newsList...)

		if len(l) < 10 {
			break
		}
	}

	filename := writeCSV(newsList, query)
	fmt.Println("검색을 종료합니다. " + strconv.Itoa(len(newsList)) + "개의 뉴스가 " + filename + " 에 저장되었습니다.")
}

func writeCSV(newsList []newsItem, query string) string {
	os.Mkdir("output", os.ModeDir)
	file, err := os.Create("output/news_" + query + "_" + time.Now().Format("2006.01.02") + ".csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	header := []string{"title", "desc", "publisher", "date", "link"}
	wErr := w.Write(header)
	checkErr(wErr)

	for _, news := range newsList {
		newsSlice := []string{news.title, news.desc, news.publisher, news.pubDate, news.link}
		nwErr := w.Write(newsSlice)
		checkErr(nwErr)
	}

	return file.Name()
}

func getPage(baseUrl string, params string, page int) []newsItem {
	var newsList []newsItem
	c := make(chan newsItem)
	params += ("&start=" + strconv.Itoa((page*10)+1))
	baseUrl += params

	res := getHttp(baseUrl)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	newsListDoc := doc.Find(".list_news")
	cards := newsListDoc.Find(".bx")
	cards.Each(func(i int, card *goquery.Selection) {
		go extractNews(card, c)
	})

	for i := 0; i < cards.Length(); i++ {
		news := <-c
		newsList = append(newsList, news)
	}

	return newsList
}

func extractNews(card *goquery.Selection, c chan<- newsItem) {
	title, _ := card.Find(".news_tit").Attr("title")
	desc := card.Find(".dsc_txt_wrap").Text()
	infos := card.Find(".info").Map(func(i int, info *goquery.Selection) string {
		remove := info.Children().Text()
		text := info.Text()
		result := strings.Replace(text, remove, "", -1)
		return result
	})
	link, _ := card.Find(".news_tit").Attr("href")

	c <- newsItem{
		title:     title,
		desc:      desc,
		publisher: infos[0],
		pubDate:   infos[2],
		link:      link,
	}
}
