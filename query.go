package bilicoin

import "net/http"

func API(url string, interceptor func(reqPoint *http.Request)) (*http.Response, error) {
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	interceptor(req)

	res, err := client.Do(req)
	return res, nil
}
