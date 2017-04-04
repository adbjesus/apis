package apis

import "fmt"
import "flag"
import "time"
import "github.com/adbjesus/apis/database"
import "github.com/adbjesus/apis/fetcher"

var trendsSearchQuery string

func init() {
	flag.StringVar(&trendsSearchQuery, "trends-search-query", "", "Trends API: Search Query")
}

func Trends() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	rows, err := database.GetResult(&db, fmt.Sprintf("select count(*) from google_trends where date_id=%s and query like '%s'", time.Now().UTC().Format("20060102"), trendsSearchQuery))
	var count int
	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	rows.Close()
	if count == 1 {
		return
	}

	var date = "today%203-m"
	body, err := fetcher.GetPage(fmt.Sprintf("http://www.google.com/trends/fetchComponent?q=%s&date=%s&cid=TIMESERIES_GRAPH_0&export=3", trendsSearchQuery, date))

	text := string(body)
	//fmt.Println(text)

	var arr []string
	t := time.Now()
	t = t.UTC()

	arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), t.Format("20060102"), trendsSearchQuery, text)
	database.InsertArray(&db, "google_trends", arr)

}
