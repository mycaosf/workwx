package workwx

import (
	"fmt"
)

type User struct {
	token
}

type UserInfo struct {
	baseResponse
	UserId  string `json:"userid"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Gender  int    `json:"gender"`
	Status  int    `json:"status"`
	QrUrl   string `json:"qr_code"`
	Address string `json:"address"`
}

func (p *User) Info(userId string) (info UserInfo, err error) {
	var url string
	if url, err = p.buildUrlInfo(userId, false); err == nil {
		if err = httpGetJson(url, nil, &info); err != nil {
			if url, err = p.buildUrlInfo(userId, true); err == nil {
				err = httpGetJson(url, nil, &info)
			}
		}
	}

	return
}

func (p *User) buildUrlInfo(userId string, force bool) (url string, err error) {
	var token string
	if token, err = p.Get(force); err == nil {
		url = fmt.Sprintf(urlUserInfo, token, userId)
	}

	return
}

const (
	urlUserInfo = "https://qyapi.weixin.qq.com/cgi-bin/user/get?access_token=%s&userid=%s"
)
