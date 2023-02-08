package workwx

import (
	"bytes"
	"encoding/json"
)

func wedrivePost(p *token, api string, req any, res any) (err error) {
	var data []byte
	if data, err = json.Marshal(req); err == nil {
		buffer := bytes.NewReader(data)

		err = p.postJson(wedriveClass, api, buffer, res)
	}

	return
}

const (
	wedriveClass = "wedrive"
)
