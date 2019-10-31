package workwx

import ()

type User struct {
	token
}

type UserInfo struct {
	baseResponse
	UserId  string `json:"userid"`
	Name    string `json:"name"`
	Mobile  string `json:"mobile"`
	Email   string `json:"email"`
	Gender  string `json:"gender"`
	Status  int    `json:"status"`
	QrUrl   string `json:"qr_code"`
	Address string `json:"address"`
}

type UserId struct {
	baseResponse
	UserId  string `json:"userid"`
	OpoenId string `json:"openid"`
}

func (p *User) Info(userId string) (info UserInfo, err error) {
	err = p.getJson(userClass, userApiGet, &info, "&userid=", userId)

	return
}

func (p *User) UserId(code string) (id UserId, err error) {
	err = p.getJson(userClass, userApiUserId, &id, "&code=", code)

	return
}

const (
	userClass     = "user"
	userApiGet    = "get"
	userApiUserId = "getuserinfo"
)
