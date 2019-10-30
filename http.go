package workwx

import (
	"encoding/json"
	"github.com/mycaosf/utils/net/httpc"
	"io"
	"io/ioutil"
	"net/http"
)

func httpGet(url string, header http.Header) (*http.Response, error) {
	return httpc.Get(url, header, nil)
}

func httpPost(url string, header http.Header, body io.Reader) (*http.Response, error) {
	return httpc.Post(url, header, body, nil)
}

func httpGetJson(url string, header http.Header, v interface{}) error {
	if resp, err := httpGet(url, header); err == nil {
		return httpParseJson(resp, v)
	} else {
		return err
	}
}

func httpPostJson(url string, header http.Header, body io.Reader, v interface{}) error {
	if resp, err := httpPost(url, header, body); err == nil {
		return httpParseJson(resp, v)
	} else {
		return err
	}
}

func httpParseJson(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if v == nil {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return json.Unmarshal(data, v)
	}

	return err
}

func httpGetBytes(url string, header http.Header) ([]byte, error) {
	if resp, err := httpGet(url, header); err == nil {
		return httpParseBytes(resp)
	} else {
		return nil, err
	}
}

func httpPostBytes(url string, header http.Header, body io.Reader) ([]byte, error) {
	if resp, err := httpPost(url, header, body); err == nil {
		return httpParseBytes(resp)
	} else {
		return nil, err
	}
}

func httpParseBytes(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
