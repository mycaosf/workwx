package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mycaosf/utils/net/httpc"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Media struct {
	token
}

func (p *Media) File(fileName string) (string, error) {
	return p.send(mediaContentTypeFile, fileName, true)
}

//only support AMR
func (p *Media) Voice(fileName string) (string, error) {
	return p.send(mediaContentTypeAmr, fileName, true)
}

//only support mp4
func (p *Media) Video(fileName string) (string, error) {
	return p.send(mediaContentTypeMp4, fileName, true)
}

func (p *Media) Jpeg(fileName string, temporary bool) (string, error) {
	return p.send(mediaContentTypeJpg, fileName, temporary)
}

func (p *Media) Png(fileName string, temporary bool) (string, error) {
	return p.send(mediaContentTypePng, fileName, temporary)
}

func (p *Media) Bmp(fileName string, temporary bool) (string, error) {
	return p.send(mediaContentTypeBmp, fileName, temporary)
}

func (p *Media) send(mediaType int, fileName string, temporary bool) (id string, err error) {
	var url string
	var fileInfo os.FileInfo

	if fileInfo, err = os.Stat(fileName); err != nil {
		return
	}
	if url, err = p.buildUrl(mediaType, temporary, false); err != nil {
		return
	}

	var buf bytes.Buffer
	var part io.Writer
	wpart := multipart.NewWriter(&buf)

	h := make(textproto.MIMEHeader)
	h.Add(httpc.HTTPHeaderContentType, mediaContentTypeStr[mediaType])
	h.Add(httpc.HTTPHeaderContentDisposition, fmt.Sprintf(mediaContentDisposition, "media", filepath.Base(fileName), fileInfo.Size()))

	var file *os.File
	if file, err = os.Open(fileName); err != nil {
		return
	}

	if part, err = wpart.CreatePart(h); err != nil {
		return
	}

	defer file.Close()
	if _, err = io.Copy(part, file); err != nil {
		wpart.Close()
		return
	}
	wpart.Close()

	httpContentType := mediaHttpContentType + wpart.Boundary()

	var body []byte
	r := bytes.NewReader(buf.Bytes())

	header := make(http.Header)
	header.Add(httpc.HTTPHeaderContentType, httpContentType)
	header.Add(httpc.HTTPHeaderContentLength, strconv.Itoa(buf.Len()))

	if body, err = httpPostBytes(url, header, r); err != nil {
		if url, err = p.buildUrl(mediaType, temporary, true); err == nil {
			r.Seek(0, 0)
			body, err = httpPostBytes(url, header, r)
		}
	}

	if err == nil {
		id, err = mediaParseId(body, temporary)
	}

	return
}

func mediaParseId(data []byte, temporary bool) (id string, err error) {
	type tempResponse struct {
		baseResponse
		MediaId string `json:"media_id"`
	}

	type persResponse struct {
		baseResponse
		Url string `json:"url"`
	}

	if temporary {
		var r tempResponse
		if err = json.Unmarshal(data, &r); err == nil {
			id = r.MediaId
		}
	} else {
		var r persResponse
		if err = json.Unmarshal(data, &r); err == nil {
			id = r.Url
		}
	}

	return
}

func (p *Media) buildUrl(mediaType int, temporary, force bool) (url string, err error) {
	var token string
	if token, err = p.Get(force); err == nil {
		if temporary {
			url = fmt.Sprintf(mediaUrlTemp, token, mediaTypeStr[mediaType])
		} else {
			url = fmt.Sprintf(mediaUrlPersistence, token)
		}

	}

	return
}

//to > 0 for range
func (p *Media) GetData(id string, from, to int) (data []byte, err error) {
	var resp *http.Response
	var header http.Header
	if to > 0 {
		header = make(http.Header)
		header.Add(httpc.HTTPHeaderRange, fmt.Sprintf("bytes=%d-%d", from, to))
	}

	if strings.HasPrefix(id, "http://") || strings.HasPrefix(id, "https://") {
		resp, err = httpGet(id, header)
	} else {
		var token string
		if token, err = p.Get(false); err == nil {
			url := fmt.Sprintf(mediaUrlGet, token, id)
			if resp, err = httpGet(url, header); err != nil {
				if tokenNew, errNew := p.Get(true); err == nil {
					if token != tokenNew {
						url := fmt.Sprintf(mediaUrlGet, tokenNew, id)
						resp, err = httpGet(url, header)
					}
				} else {
					err = errNew
				}
			}
		}
	}

	if resp != nil {
		defer resp.Body.Close()
		data, err = ioutil.ReadAll(resp.Body)
	}

	return
}

var (
	mediaContentTypeStr = []string{
		"application/octet-stream",
		"image/jpg",
		"image/png",
		"image/bmp",
		"voice/amr",
		"voice/mp4",
	}
	mediaTypeStr = []string{
		"file", "image", "image", "image", "voice", "video",
	}
)

const (
	mediaContentTypeFile = iota
	mediaContentTypeJpg
	mediaContentTypePng
	mediaContentTypeBmp
	mediaContentTypeAmr
	mediaContentTypeMp4

	mediaHttpContentType    = "multipart/form-data; boundary="
	mediaContentDisposition = "form-data; name=\"%s\";filename=\"%s\"; filelength=%v"
	mediaUrlTemp            = "https://qyapi.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s"
	mediaUrlPersistence     = "https://qyapi.weixin.qq.com/cgi-bin/media/uploadimg?access_token=%s"
	mediaUrlGet             = "https://qyapi.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s"
)
