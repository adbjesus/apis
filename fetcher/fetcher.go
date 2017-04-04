package fetcher

import "net/http"
import "io/ioutil"

func GetPage(st string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", st, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 5.1; rv:31.0) Gecko/20100101 Firefox/31.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	//resp, err := http.Get(st)
	//if err != nil {
	//	return nil, err
	//}
	t, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return t, nil
}
