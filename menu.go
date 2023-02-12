package workwx

import (
	"bytes"
	"strconv"
)

type Menu struct {
	Token
	agentId int
}

func (p *Menu) SetAgentId(agentId int) {
	p.agentId = agentId
}

func (p *Menu) Create(str string) (ret Error) {
	data := bytes.NewReader([]byte(str))
	if err := p.PostJson(menuClass, menuApiCreate, data, &ret, agentIdStr, strconv.Itoa(p.agentId)); err != nil {
		ret.SetError(err)
	}

	return
}

func (p *Menu) Data() (string, error) {
	if data, err := p.GetBytes(menuClass, menuApiGet, nil, agentIdStr, strconv.Itoa(p.agentId)); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func (p *Menu) Delete() (ret Error) {
	p.getJson(menuClass, menuApiDelete, &ret, agentIdStr, strconv.Itoa(p.agentId))

	return
}

const (
	menuClass     = "menu"
	menuApiCreate = "create"
	menuApiGet    = "get"
	menuApiDelete = "delete"
	agentIdStr    = "&agentid="
)
