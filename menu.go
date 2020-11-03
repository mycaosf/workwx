package workwx

import (
	"bytes"
	"strconv"
)

type Menu struct {
	token
	agentId int
}

func (p *Menu) SetAgentId(agentId int) {
	p.agentId = agentId
}

func (p *Menu) Create(str string) error {
	data := bytes.NewReader([]byte(str))
	var e Error
	if err := p.postJson(menuClass, menuApiCreate, data, &e, agentIdStr, strconv.Itoa(p.agentId)); err == nil {
		return e.parse()
	} else {
		return err
	}
}

func (p *Menu) Data() (string, error) {
	if data, err := p.getBytes(menuClass, menuApiGet, nil, agentIdStr, strconv.Itoa(p.agentId)); err == nil {
		return string(data), nil
	} else {
		return "", err
	}
}

func (p *Menu) Delete() error {
	var e Error
	if err := p.getJson(menuClass, menuApiDelete, &e, agentIdStr, strconv.Itoa(p.agentId)); err == nil {
		return e.parse()
	} else {
		return err
	}
}

const (
	menuClass     = "menu"
	menuApiCreate = "create"
	menuApiGet    = "get"
	menuApiDelete = "delete"
	agentIdStr    = "&agentid="
)
