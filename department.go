package workwx

import (
	"strconv"
)

type Department struct {
	Token
}

type DepartmentItem struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	ParentId int    `json:"parentid"`
	Order    int    `json:"order"`
}

type DepartmentListResponse struct {
	Error
	Items []DepartmentItem `json:"department"`
}

// id < 0 if for all
func (p *Department) List(id int) (ret DepartmentListResponse) {
	if id > 0 {
		p.getJson(departmentClass, departmentApiList, &ret, "&id=", strconv.Itoa(id))
	} else {
		p.getJson(departmentClass, departmentApiList, &ret)
	}

	return
}

const (
	departmentClass   = "department"
	departmentApiList = "list"
	departmentIdStr   = "&department_id="
	fetchChildStr     = "&fetch_child="
)
