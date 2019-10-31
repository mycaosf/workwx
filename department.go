package workwx

import (
	"strconv"
)

type DepartmentItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentid"`
	Order    int    `json:"order"`
}

type DepartmentList struct {
	baseResponse
	Items []DepartmentItem `json:"department"`
}

type Department struct {
	token
}

//id < 0 if for all
func (p *Department) List(id int) (ret DepartmentList, err error) {
	if id > 0 {
		err = p.getJson(departmentClass, departmentApiList, &ret, "&id=", strconv.Itoa(id))
	} else {
		err = p.getJson(departmentClass, departmentApiList, &ret)
	}

	return
}

const (
	departmentClass   = "department"
	departmentApiList = "list"
)
