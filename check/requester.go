package check

import (
	"bytes"
	_ "fmt"
	"io"
	_ "io/ioutil"
	"net/http"
)

type Requester interface {
	DoRequest(url string, method string) (int, error) // do simple HTTP request with method (GET, POST etc.) and URL
}

type HttpRequester struct {
	client *http.Client
}

func NewHttpRequester() *HttpRequester {
	return &HttpRequester{
		client: &http.Client{},
	}
}

func (httpRequester *HttpRequester) DoRequest(url string, method string) (int, error) {
	var jsonStr = []byte(``)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return 500, err
	}
	resp, err := httpRequester.client.Do(req)
	if err != nil {
		return 500, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	return resp.StatusCode, nil
}
