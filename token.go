package workwx

import (
	"errors"
	"fmt"
	"github.com/mycaosf/utils/net/httpc"
	"io"
	"net/http"
	"time"
)

type token struct {
	corpId string
	secret string
	token  string
	expire time.Time
}

type baseResponse struct {
	ErrCode int    `json:"errorcode"`
	ErrMsg  string `json:"errmsg"`
}

type accessTokenGetResponse struct {
	baseResponse
	Token  string `json:"access_token"`
	Expire int    `json:"expires_in"`
}

type Error struct {
	baseResponse
}

func (p *baseResponse) Error() string {
	return fmt.Sprintf("errcode: %d, errmsg: %s", p.ErrCode, p.ErrMsg)
}

func (p *baseResponse) parse() error {
	if p.ErrCode == 0 {
		return nil
	} else {
		return errors.New(p.Error())
	}
}

//return token
func (p *token) Get(force bool) (token string, err error) {
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

func (p *token) get() error {
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

func (p *token) urlToken(class, api string, force bool, exts ...string) (url string, err error) {
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

func (p *token) getJson(class, api string, v interface{}, exts ...string) (err error) {
	var url string
	if url, err = p.urlToken(class, api, false, exts...); err == nil {
		if err = httpGetJson(url, nil, v); err != nil {
			if url, err = p.urlToken(class, api, true, exts...); err == nil {
				err = httpGetJson(url, nil, v)
			}
		}
	}

	return
}

func (p *token) getBytes(class, api string, header http.Header, exts ...string) (data []byte, err error) {
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

func (p *token) postJson(class, api string, data io.ReadSeeker, r interface{}, exts ...string) (err error) {
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

//set cropId and secret
func (p *token) Set(corpId, secret string) {
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
