package spider

import (
	"io/ioutil"
	"net/http"
)

func Spider(url string) (content string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return string(body), err
}
