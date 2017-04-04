package apis

import "fmt"
import "flag"
import "time"
import "strings"

import "github.com/adbjesus/apis/parser"
import "github.com/adbjesus/apis/database"
import "github.com/adbjesus/apis/fetcher"

var fbPage string
var fbApiId string
var fbApiSecret string

func init() {
	flag.StringVar(&fbPage, "fb-page", "", "Facebook API: Page to get info on")
	flag.StringVar(&fbApiId, "fb-api-id", "", "Faceebook API: Api's id")
	flag.StringVar(&fbApiSecret, "fb-api-secret", "", "Facebook API: Api's secret")
}

func Facebook() {
	//var db database.DB
	flag.Parse()

	err := database.Connect(&db)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}
	access_token := getAccessToken()
	if access_token == "" {
		return
	}
	getInfo(fbPage, access_token)
	getFeed(fbPage, access_token)
}

func getAccessToken() string {
	loc := fmt.Sprintf("https://graph.facebook.com/oauth/access_token?client_id=%s&client_secret=%s&grant_type=client_credentials", fbApiId, fbApiSecret)
	page, err := fetcher.GetPage(loc)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return ""
	}
	token := string(page[:])
	ind := strings.Index(token, "access_token=")
	if ind == -1 {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println("Error getting access token")
		return ""
	}
	return token[13:]
}

func mapToArray(m map[string]interface{}, keys []string) []string {
	var arr []string

	for _, k := range keys {
		v, ok := m[k]
		if ok {
			switch v.(type) {
			case string:
				arr = append(arr, v.(string))
				break
			case float64:
				arr = append(arr, fmt.Sprintf("%d", int(v.(float64))))
				break
			case bool:
				arr = append(arr, fmt.Sprintf("%t", v.(bool)))
				break
			default:
				arr = append(arr, "")
				break
			}
		} else {
			arr = append(arr, "")
		}
	}

	return arr
}

func getInfo(page string, access_token string) {
	var info map[string]interface{}
	var err error
	var url string

	url = fmt.Sprintf("https://graph.facebook.com/%s?access_token=%s", page, access_token)
	err = parser.ParseWebJSON(url, &info)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return
	}

	t := time.Now()
	t = t.UTC()
	var arr []string
	var keys = []string{"name", "likes", "talking_about_count"}

	arr = append(arr, t.Format("2006/01/02 15:04:05 UTC"), t.Format("20060102"), page)
	arr = append(arr, mapToArray(info, keys)...)
	database.InsertArray(&db, "fb_info", arr)
}

func getPostCount(post string, access_token string, what string) int {
	type Likes struct {
		Summary map[string]interface{} `json:"summary"`
	}

	var url string
	var err error
	var data Likes

	url = fmt.Sprintf("https://graph.facebook.com/%s/%s?summary=1&access_token=%s", post, what, access_token)
	err = parser.ParseWebJSON(url, &data)
	if err != nil {
		fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
		fmt.Println(err)
		return 0
	}

	if data.Summary != nil && data.Summary["total_count"] != nil {
		return int(data.Summary["total_count"].(float64))
	}

	return 0
}

func getFeed(page string, access_token string) {
	type PostInfo struct {
		Id          string                 `json:"id"`
		CreatedTime string                 `json:"created_time"`
		Message     string                 `json:"message"`
		Type        string                 `json:"type"`
		Shares      map[string]interface{} `json:"shares"`
	}
	type Data struct {
		Data   []PostInfo             `json:"data"`
		Paging map[string]interface{} `json:"paging"`
	}

	var err error
	var url string
	var count = 0

	url = fmt.Sprintf("https://graph.facebook.com/%s/feed?access_token=%s", page, access_token)

	for err == nil {
		var data Data
		err = parser.ParseWebJSON(url, &data)
		if err != nil {
			fmt.Print(time.Now().Format("2006/01/02 15:04:05 : "))
			fmt.Println(err)
			return
		}
		for _, post := range data.Data {
			likes_count := getPostCount(post.Id, access_token, "likes")
			comments_count := getPostCount(post.Id, access_token, "comments")
			var arr []string
			arr = append(arr, post.CreatedTime, page, post.Id, post.Message, post.Type, fmt.Sprintf("%d", likes_count), fmt.Sprintf("%d", comments_count))
			if post.Shares == nil || post.Shares["count"] == nil {
				arr = append(arr, "0")
			} else {
				arr = append(arr, fmt.Sprintf("%d", int(post.Shares["count"].(float64))))
			}
			err = database.InsertArray(&db, "fb_posts", arr)
			if err != nil {
				break
			}
		}
		count += len(data.Data)
		if data.Paging == nil || data.Paging["next"] == nil {
			break
		}
		url = data.Paging["next"].(string)
	}
}
