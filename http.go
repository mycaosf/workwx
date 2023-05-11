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

func httpPostBytes(url string, header http.Header, body []byte) ([]byte, error) {
	c := httpc.Client{
		Header: header,
	}

	return c.PostBytes(url, body)
}

func httpPostJson(url string, header http.Header, body, v interface{}) error {
	c := httpc.Client{
		Header: header,
	}

	return c.PostJSON(url, body, v)
}
