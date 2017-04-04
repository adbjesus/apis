package parser

import "github.com/adbjesus/apis/fetcher"
import "encoding/xml"
import "encoding/json"
import "strings"
import "log"

func ParseXML(data []byte, v interface{}) error {

	err := xml.Unmarshal(data, v)

	return err
}

func ParseWebXML(page string, v interface{}) error {

	data, err := fetcher.GetPage(page)

	if err != nil {
		return err
	}

	err = ParseXML(data, v)

	return err
}

func ParseJSON(data []byte, v interface{}) error {

	err := json.Unmarshal(data, v)

	return err
}

func ParseWebJSON(page string, v interface{}) error {

	data, err := fetcher.GetPage(page)

	if err != nil {
		return err
	}
	err = ParseJSON(data, v)

	return err
}

func ParseText(text, ini, end string) string {
	if strings.Index(text, ini) == -1 {
		log.Println("Could not find beggining")
		return ""
	} else {
		text = text[strings.Index(text, ini)+len(ini):]
	}

	if strings.Index(text, end) == -1 {
		log.Println("Could not find ending")
		return ""
	} else {
		text = text[:strings.Index(text, end)]
	}
	return text

}
