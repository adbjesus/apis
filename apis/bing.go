package apis

import "fmt"
import "flag"
import "time"
import "strings"
import "github.com/adbjesus/apis/database"
import "github.com/adbjesus/apis/fetcher"
import "github.com/adbjesus/apis/parser"

var bingSearchQuery string

func init() {
	flag.StringVar(&bingSearchQuery, "bing-search-query", "", "Bing API: Search Query")
}

func Bing() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	body, err := fetcher.GetPage(fmt.Sprintf("http://bing.com/search?q=%s", bingSearchQuery))

	var count string
	count = parser.ParseText(string(body), "<span class=\"sb_count\">", " resultados</span>")
	if count == "" {
		count = "0"
	}

	count = strings.Replace(count, "&#160;", "", -1)

	var arr []string
	t := time.Now()
	t = t.UTC()

	arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), t.Format("20060102"), bingSearchQuery, count)
	database.InsertArray(&db, "bing_search", arr)
}
