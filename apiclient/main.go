package main

//import "fmt"
import "flag"
import "github.com/adbjesus/apis/apis"

var apiFlag string

func init() {
	flag.StringVar(&apiFlag, "api", "nextbus", "API we are using")
}

func main() {
	flag.Parse()

	switch apiFlag {
	case "nextbus":
		apis.Nextbus()
		break
	/*case "here":
	apis.Here()
	break*/
	case "facebook":
		apis.Facebook()
		break
	case "twitter":
		apis.Twitter()
		break
	case "google":
		apis.Google()
		break
	case "bing":
		apis.Bing()
		break
	case "news":
		apis.News()
		break
	case "trends":
		apis.Trends()
		break
	case "first-page":
		apis.GoogleFirstPage()
		break
	default:
		flag.Usage()
		return
	}
}
