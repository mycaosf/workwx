package workwx

import (
	"github.com/mycaosf/utils/net/httpc"
	"io"
	"net/http"
)

func httpGet(url string, header http.Header) (*http.Response, error) {
	return httpc.Get(url, header, nil)
}

func httpPost(url string, header http.Header, body io.Reader) (*http.Response, error) {
	return httpc.Post(url, header, body, nil)
}
