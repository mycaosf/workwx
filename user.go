package workwx

import (
	"strconv"
)

type User struct {
	token
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

func (p *User) Info(userId string) (info UserInfo, err error) {
	type UserInfoResponse struct {
		baseResponse
		UserInfo
	}

	var r UserInfoResponse
	if err = p.getJson(userClass, userApiGet, &r, "&userid=", userId); err == nil {
		info = r.UserInfo
	}

	return
}

func (p *User) UserId(code string) (id UserId, err error) {
	type UserIdResponse struct {
		baseResponse
		UserId
	}

	var r UserIdResponse
	if err = p.getJson(userClass, userApiUserId, &r, "&code=", code); err == nil {
		id = r.UserId
	}

	return
}

func (p *User) DepartmentUsers(departmentId int, fetchChild bool) (users []UserInfoBase, err error) {
	type DepartmentUsersResponse struct {
		baseResponse
		Users []UserInfoBase `json:"userlist"`
	}

	var r DepartmentUsersResponse
	department := strconv.Itoa(departmentId)
	fetch := "0"
	if fetchChild {
		fetch = "1"
	}

	if err = p.getJson(userClass, userApiSimpleList, &r, departmentIdStr, department, fetchChildStr, fetch); err == nil {
		users = r.Users
	}

	return
}

func (p *User) DepartmentUsersDetail(departmentId int, fetchChild bool) (users []UserInfo, err error) {
	type DepartmentUsersResponse struct {
		baseResponse
		Users []UserInfo `json:"userlist"`
	}

	var r DepartmentUsersResponse
	department := strconv.Itoa(departmentId)
	fetch := "0"
	if fetchChild {
		fetch = "1"
	}

	if err = p.getJson(userClass, userApiList, &r, departmentIdStr, department, fetchChildStr, fetch); err == nil {
		users = r.Users
	}

	return
}

const (
	userClass         = "user"
	userApiGet        = "get"
	userApiUserId     = "getuserinfo"
	userApiSimpleList = "simplelist"
	userApiList       = "list"
)
