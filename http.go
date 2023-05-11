package workwx

import (
	"github.com/mycaosf/httpc"
	"io"
	"net/http"
)

func httpGet(url string, header http.Header) (*http.Response, error) {
	c := httpc.Client{
		Header: header,
	}

	return c.Get(url)
}

func httpGetBytes(url string, header http.Header) ([]byte, error) {
	c := httpc.Client{
		Header: header,
	}

	return c.GetBytes(url)
}

func httpGetJson(url string, header http.Header, v interface{}) error {
	c := httpc.Client{
		Header: header,
	}

	return c.GetJSON(url, v)
}

func httpPost(url string, header http.Header, body io.Reader) (*http.Response, error) {
	c := httpc.Client{
		Header: header,
	}

	return c.Post(url, body)
}

func httpPostBytes(url string, header http.Header, body io.Reader) ([]byte, error) {
	if resp, err := httpPost(url, header, body); err == nil {
		defer resp.Body.Close()

		return io.ReadAll(resp.Body)
	} else {
		return nil, err
	}
}

func httpPostJson(url string, header http.Header, body io.Reader, v interface{}) error {
	c := httpc.Client{
		Header: header,
	}

	return c.PostJSON(url, body, v)
}
