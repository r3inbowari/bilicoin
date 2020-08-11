package bilicoin

import (
	"net/http"
	"strings"
)

func GET(url string, interceptor func(reqPoint *http.Request)) (*http.Response, error) {
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, nil
}

func Post2(url string, interceptor func(reqPoint *http.Request), body string) (*http.Response, error) {
	method := "POST"

	client := &http.Client{}
	payload := strings.NewReader(body)
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, err
}

func Post(url string, interceptor func(reqPoint *http.Request)) (*http.Response, error) {
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if interceptor != nil {
		interceptor(req)
	}
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	res, err := client.Do(req)
	return res, err
}
