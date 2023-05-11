package workwx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mycaosf/httpc"
	"io"
	"net/http"
	"time"
)

type Token struct {
	corpId string
	secret string
	token  string
	expire time.Time
}

type accessTokenGetResponse struct {
	Error
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

// return token string
func (p *Token) Get(force bool) (token string, err error) {
	if p.corpId == "" {
		err = errTokenNotInit

		return
	} else if p.token == "" || force {
		if err = p.get(); err != nil {
			return
		}
	} else {
		t := time.Now()
		if t.After(p.expire) {
			if err = p.get(); err != nil {
				return
			}
		}
	}

	token = p.token

	return
}

func (p *Token) get() error {
	var ret accessTokenGetResponse
	p.token = ""
	url := urlBase + "gettoken?corpid=" + p.corpId + "&corpsecret=" + p.secret

	err := httpGetJson(url, nil, &ret)
	if err == nil {
		p.token = ret.Token
		p.expire = time.Now().Add(time.Duration(ret.Expire) * time.Second)
	}

	return err
}

func (p *Token) urlToken(class, api string, force bool, exts ...string) (url string, err error) {
	var token string
	if token, err = p.Get(force); err == nil {
		url = urlBase + fmt.Sprintf("%s/%s?access_token=%s", class, api, token)
		if count := len(exts); count > 0 {
			for i := 0; i < count; i++ {
				url += exts[i]
			}
		}
	}

	return
}

func (p *Token) getJson(class, api string, res any, exts ...string) {
	if err := p.GetJson(class, api, res, exts...); err != nil {
		if e, ok := res.(IError); ok {
			e.SetError(err)
		}
	}
}

func (p *Token) GetJson(class, api string, res any, exts ...string) (err error) {
	var url string
	if url, err = p.urlToken(class, api, false, exts...); err == nil {
		if err = httpGetJson(url, nil, res); err != nil {
			if url, err = p.urlToken(class, api, true, exts...); err == nil {
				err = httpGetJson(url, nil, res)
			}
		}
	}

	return
}

func (p *Token) GetBytes(class, api string, header http.Header, exts ...string) (data []byte, err error) {
	var url string
	if url, err = p.urlToken(class, api, false, exts...); err == nil {
		if data, err = httpGetBytes(url, header); err != nil {
			if url, err = p.urlToken(class, api, true, exts...); err == nil {
				data, err = httpGetBytes(url, header)
			}
		}
	}

	return
}

func (p *Token) postJson(class, api string, req, res any, exts ...string) {
	data, err := json.Marshal(req)
	if err == nil {
		buffer := bytes.NewReader(data)
		err = p.PostJson(class, api, buffer, res)
	}

	if err != nil {
		if e, ok := res.(IError); ok {
			e.SetError(err)
		}
	}
}

func (p *Token) PostJson(class, api string, data io.ReadSeeker, r any, exts ...string) (err error) {
	var url string
	if url, err = p.urlToken(class, api, false, exts...); err == nil {
		header := make(http.Header)
		header.Add(httpc.HTTPHeaderContentType, contentJson)

		if err = httpPostJson(url, header, data, r); err != nil {
			data.Seek(0, 0)
			if url, err = p.urlToken(class, api, true, exts...); err == nil {
				err = httpPostJson(url, header, data, r)
			}
		}
	}

	return
}

// set cropId and secret
func (p *Token) Set(corpId, secret string) {
	p.corpId = corpId
	p.secret = secret
	p.token = ""
}

var (
	errTokenNotInit = errors.New("token isn't initialized.")
)

const (
	urlBase     = "https://qyapi.weixin.qq.com/cgi-bin/"
	contentJson = "application/json"
)
