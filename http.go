package workwx

import (
	"github.com/mycaosf/utils/net/httpc"
	"io"
	"net/http"
)

func httpGet(url string, header http.Header) (*http.Response, error) {
	return httpc.Get(url, header, nil)
}

func httpGetBytes(url string, header http.Header) ([]byte, error) {
	return httpc.GetBytes(url, header, nil)
}

func httpGetJson(url string, header http.Header, v interface{}) error {
	return httpc.GetJson(url, header, nil, v)
}

func httpPost(url string, header http.Header, body io.Reader) (*http.Response, error) {
	return httpc.Post(url, header, body, nil)
}

func httpPostBytes(url string, header http.Header, body io.Reader) ([]byte, error) {
	return httpc.PostBytes(url, header, body, nil)
}

func httpPostJson(url string, header http.Header, body io.Reader, v interface{}) error {
	return httpc.PostJson(url, header, body, nil, v)
}
