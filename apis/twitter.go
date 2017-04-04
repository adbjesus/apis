package apis

import "fmt"
import "flag"
import "time"
import "net/url"
import "github.com/ChimeraCoder/anaconda"
import "github.com/adbjesus/apis/database"

var twitterQuery string
var twitterScreenName string
var twitterSearchBool bool
var twitterTimelineBool bool
var twitterApiKey string
var twitterApiSecret string
var twitterAccessKey string
var twitterAccessSecret string

func init() {
	flag.StringVar(&twitterQuery, "twitter-query", "", "Twitter API: Search query")
	flag.StringVar(&twitterScreenName, "twitter-screenname", "", "Twitter API: Page to get info on")
	flag.BoolVar(&twitterSearchBool, "twitter-search", false, "Twitter API: Boolean, do search?")
	flag.BoolVar(&twitterTimelineBool, "twitter-timeline", false, "Twitter API: Boolean, do timeline?")
	flag.StringVar(&twitterApiKey, "twitter-api-key", "", "Twitter API: Api's key")
	flag.StringVar(&twitterApiSecret, "twitter-api-secret", "", "Twitter API: Api's secret")
	flag.StringVar(&twitterAccessKey, "twitter-access-key", "", "Twitter API: Access token's key")
	flag.StringVar(&twitterAccessSecret, "twitter-access-secret", "", "Twitter API: Access token's secret")
}

func twitterSearch(api *anaconda.TwitterApi) {
	v := url.Values{}
	v.Set("count", "100")
	var last_id int64
	var count int
	for count = 0; ; count++ {
		var result []anaconda.Tweet
		result, err := api.GetSearch(twitterQuery, v)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(result) == 0 {
			break
		}

		for _, tweet := range result {
			last_id = tweet.Id
			if len(tweet.Text) >= 3 && tweet.Text[:3] == "RT " {
				continue
			}
			var arr []string
			arr = append(arr, tweet.CreatedAt, twitterQuery, tweet.IdStr, tweet.Text, fmt.Sprintf("%d", tweet.RetweetCount), fmt.Sprintf("%d", tweet.FavoriteCount))
			err := database.InsertArray(&db, "twitter_search", arr)
			if err != nil {
				return
			}
		}
		v.Set("max_id", fmt.Sprintf("%d", last_id-1))
	}
}

func twitterTimeline(api *anaconda.TwitterApi) {
	v := url.Values{}
	v.Set("count", "100")
	v.Set("used_id", twitterScreenName)
	v.Set("screen_name", twitterScreenName)
	var last_id int64
	var count int
	for count = 0; ; count++ {
		var result []anaconda.Tweet
		result, err := api.GetUserTimeline(v)
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(result) == 0 {
			break
		}

		for _, tweet := range result {
			last_id = tweet.Id
			if len(tweet.Text) >= 3 && tweet.Text[:3] == "RT " {
				continue
			}
			var arr []string
			arr = append(arr, tweet.CreatedAt, twitterScreenName, tweet.IdStr, tweet.Text, fmt.Sprintf("%d", tweet.RetweetCount), fmt.Sprintf("%d", tweet.FavoriteCount))
			err := database.InsertArray(&db, "twitter_timeline", arr)
			if err != nil {
				return
			}

		}
		v.Set("max_id", fmt.Sprintf("%d", last_id-1))
	}
}

func Twitter() {
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	anaconda.SetConsumerKey(twitterApiKey)
	anaconda.SetConsumerSecret(twitterApiSecret)

	api := anaconda.NewTwitterApi(twitterAccessKey, twitterAccessSecret)
	api.ReturnRateLimitError(true)
	api.EnableThrottling(6*time.Second, 5)

	if twitterSearchBool {
		twitterSearch(api)
	}
	if twitterTimelineBool {
		twitterTimeline(api)
	}

}
