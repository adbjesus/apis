package apis

import "fmt"
import "log"
import "flag"
import "time"
import "strings"
import "code.google.com/p/go.net/html"
import "github.com/adbjesus/apis/database"
import "github.com/adbjesus/apis/fetcher"
import "github.com/adbjesus/apis/parser"

var googleSearchQuery string

func init() {
	flag.StringVar(&googleSearchQuery, "google-search-query", "", "Google API: Search Query")
}

func Google() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := fetcher.GetPage(fmt.Sprintf("http://google.com/search?q=%s&hl=en", googleSearchQuery))
	if err != nil {
		log.Println(err)
		return
	}

	var count string
	count = parser.ParseText(string(body), "id=\"resultStats\">About ", " results</div>")
	if count == "" {
		count = "0"
	}

	var arr []string
	t := time.Now()
	t = t.UTC()

	arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), t.Format("20060102"), googleSearchQuery, strings.Replace(count, ",", "", -1))
	database.InsertArray(&db, "google_search", arr)
}

func GoogleFirstPage() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		log.Println(err)
		return
	}

	var done = false
	var body []byte

	for i := 0; i < 5 && !done; i++ {
		body, err = fetcher.GetPage(fmt.Sprintf("http://google.com/search?q=%s&tbs=qdr:d&start=%d", googleSearchQuery, 10*i))
		//fmt.Println(string(body))
		links := parser.ParseText(string(body), "<div id=\"ires\">", "</html>") + "</html>"

		for !done {
			links = parser.ParseText(links, "<li class=\"g\">", "</html>") + "</html>"
			if links == "</html>" {
				break
			}
			link := parser.ParseText(links, "<a href=\"/url?q=", "\">")
			link = html.UnescapeString(link)
			link = strings.Replace(link, "%3F", "?", -1)
			link = strings.Replace(link, "%3D", "=", -1)
			link = link[:strings.Index(link, "&sa=")]
			body, err = fetcher.GetPage(link)

			if link != "" {
				var arr []string
				t := time.Now()
				t = t.UTC()

				arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), googleSearchQuery, link, string(body))
				err := database.InsertArray(&db, "google_first_page", arr)
				if err != nil {
					continue
				} else {
					done = true
				}

			}
		}
	}

}
