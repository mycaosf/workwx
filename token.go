package workwx

import (
	"errors"
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
	url := urlGetToken + "corpid=" + p.corpId + "&corpsecret=" + p.secret

	err := httpGetJson(url, nil, &ret)
	if err == nil {
		p.token = ret.Token
		p.expire = time.Now().Add(time.Duration(ret.Expire) * time.Second)
	}

	return err
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
	urlGetToken = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?"
)
