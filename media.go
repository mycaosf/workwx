package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mycaosf/utils/net/httpc"
	"io"
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

// only support AMR
func (p *Media) Voice(fileName string) (string, error) {
	return p.send(mediaContentTypeAmr, fileName, true)
}

// only support mp4
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
	if temporary {
		url, err = p.urlToken(mediaClass, mediaApiUpload, force, "&type=", mediaTypeStr[mediaType])
	} else {
		url, err = p.urlToken(mediaClass, mediaApiUploadImg, force)
	}

	return
}

// fromTo[0]: from, fromTo[1]: to
func (p *Media) GetData(id string, fromTo ...int) (data []byte, err error) {
	var header http.Header
	if len(fromTo) > 1 {
		header = make(http.Header)
		header.Add(httpc.HTTPHeaderRange, fmt.Sprintf("bytes=%d-%d", fromTo[0], fromTo[1]))
	}

	if strings.HasPrefix(id, "http://") || strings.HasPrefix(id, "https://") {
		data, err = httpGetBytes(id, header)
	} else {
		data, err = p.getBytes(mediaClass, mediaApiGet, header, "&media_id=", id)
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
	mediaClass              = "media"
	mediaApiUpload          = "upload"
	mediaApiUploadImg       = "uploadimg"
	mediaApiGet             = "get"
)
