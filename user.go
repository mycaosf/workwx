package workwx

import (
	"strconv"
)

type User struct {
	Token
}

type UserInfoBase struct {
	UserId      string `json:"userid"`
	Name        string `json:"name"`
	Departments []int  `json:"department"`
	OpenUserId  string `json:"open_userid"`
}

type UserInfo struct {
	UserInfoBase
	Order          []int  `json:"order"`
	Leader         []int  `json:"is_leader_in_dept"`
	Mobile         string `json:"mobile"`
	Telephone      string `json:"telephone"`
	Position       string `json:"position"`
	Gender         string `json:"gender"`
	Email          string `json:"email"`
	Status         int    `json:"status"`
	QrUrl          string `json:"qr_code"`
	Avatar         string `json:"avatar"`
	DepartmentMain int    `json:"main_department"`
	NameEnglish    string `json:"english_name"`
	Address        string `json:"address"`
}

type UserId struct {
	UserId  string `json:"userid"`
	OpoenId string `json:"openid"`
}

type UserInfoResponse struct {
	Error
	UserInfo
}

type UserIdResponse struct {
	Error
	UserId
}

type DepartmentUsersResponse struct {
	Error
	Users []UserInfoBase `json:"userlist"`
}

type DepartmentUsersDetailResponse DepartmentUsersResponse

func (p *User) Info(userId string) (ret UserInfoResponse) {
	p.getJson(userClass, userApiGet, &ret, "&userid=", userId)

	return
}

func (p *User) UserId(code string) (ret UserIdResponse) {
	p.getJson(userClass, userApiUserId, &ret, "&code=", code)

	return
}

func (p *User) DepartmentUsers(departmentId int, fetchChild bool) (ret DepartmentUsersResponse) {
	department := strconv.Itoa(departmentId)
	fetch := "0"
	if fetchChild {
		fetch = "1"
	}

	p.getJson(userClass, userApiSimpleList, &ret, departmentIdStr, department, fetchChildStr, fetch)

	return
}

func (p *User) DepartmentUsersDetail(departmentId int, fetchChild bool) (ret DepartmentUsersDetailResponse) {
	department := strconv.Itoa(departmentId)
	fetch := "0"
	if fetchChild {
		fetch = "1"
	}

	p.getJson(userClass, userApiList, &ret, departmentIdStr, department, fetchChildStr, fetch)

	return
}

const (
	userClass         = "user"
	userApiGet        = "get"
	userApiUserId     = "getuserinfo"
	userApiSimpleList = "simplelist"
	userApiList       = "list"
)
