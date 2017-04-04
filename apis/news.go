package apis

import "fmt"
import "flag"
import "time"
import "github.com/adbjesus/apis/database"
import "github.com/adbjesus/apis/fetcher"

var newsSearchQuery string

func init() {
	flag.StringVar(&newsSearchQuery, "news-search-query", "", "Google News API: Search Query")
}

func News() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	body, err := fetcher.GetPage(fmt.Sprintf("https://ajax.googleapis.com/ajax/services/search/news?v=1.0&q=%s", newsSearchQuery))

	var arr []string
	t := time.Now()
	t = t.UTC()

	arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), t.Format("20060102"), newsSearchQuery, string(body))
	database.InsertArray(&db, "google_news", arr)
}
