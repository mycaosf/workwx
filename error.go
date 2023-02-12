package workwx

import (
	"fmt"
)

type Error struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	err     error
}

type IError interface {
	error
	SetError(error)
	GetError() error
	IsWeworkError() bool
}

func (p *Error) Error() string {
	if p.err != nil {
		return p.err.Error()
	} else if p.ErrCode != 0 {
		return fmt.Sprintf("errcode: %d, errmsg: %s", p.ErrCode, p.ErrMsg)
	} else {
		return ""
	}
}

func (p *Error) SetError(err error) {
	p.err = err
}

func (p *Error) GetError() error {
	if p.err != nil {
		return p.err
	} else if p.ErrCode == 0 {
		return nil
	} else {
		return p
	}
}

func (p *Error) IsWeworkError() bool {
	return p.err == nil && p.ErrCode != 0
}
