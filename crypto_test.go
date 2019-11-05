package workwx

import (
	"bytes"
	"encoding/base64"
	"testing"
)

func newCrypto(t *testing.T) (*Crypto, error) {
	enc := base64.StdEncoding
	if key, err := enc.DecodeString(encodingAesKeyTest + "="); err != nil {
		t.Error("aesKey decode error:", err)
		return nil, err
	} else if c, err := NewCrypto(tokenTest, corpIdTest, key); err != nil {
		t.Error("NewCrypto failed:", err)

		return nil, err
	} else {
		return c, nil
	}
}

func TestSignature(t *testing.T) {
	if c, err := newCrypto(t); err == nil {
		signature := c.Signature(timestampsTest, nonceTest, msgTest)
		if signature != signatureTest {
			t.Error("signature is not correct")
		}
	}
}

func TestDecrypt(t *testing.T) {
	if c, err := newCrypto(t); err == nil {
		data, err := c.Decrypt(msgTest)
		if err != nil {
			t.Error("Decrypt failed:", err)
		} else if bytes.Compare(data, []byte(msgData)) != 0 {
			t.Error("Decrypt data error")
		}
	}
}

func TestEncrypt(t *testing.T) {
	if c, err := newCrypto(t); err == nil {
		msgStr := c.Encrypt([]byte(msgData))
		data, err := c.Decrypt(msgStr)

		if err != nil {
			t.Error("Decrypt failed:", err)
		} else if bytes.Compare(data, []byte(msgData)) != 0 {
			t.Error("Decrypt data error")
		}
	}
}

const (
	corpIdTest         = "wx5823bf96d3bd56c7"
	tokenTest          = "QDG6eK"
	encodingAesKeyTest = "jWmYm7qr5nMoAUwZRjGtBxmz3KA1tkAj3ykkR6q2B2C"

	signatureTest  = "477715d11cdb4164915debcba66cb864d751f3e6"
	timestampsTest = "1409659813"
	nonceTest      = "1372623149"
	msgTest        = "RypEvHKD8QQKFhvQ6QleEB4J58tiPdvo+rtK1I9qca6aM/wvqnLSV5zEPeusUiX5L5X/0lWfrf0QADHHhGd3QczcdCUpj911L3vg3W/sYYvuJTs3TUUkSUXxaccAS0qhxchrRYt66wiSpGLYL42aM6A8dTT+6k4aSknmPj48kzJs8qLjvd4Xgpue06DOdnLxAUHzM6+kDZ+HMZfJYuR+LtwGc2hgf5gsijff0ekUNXZiqATP7PF5mZxZ3Izoun1s4zG4LUMnvw2r+KqCKIw+3IQH03v+BCA9nMELNqbSf6tiWSrXJB3LAVGUcallcrw8V2t9EL4EhzJWrQUax5wLVMNS0+rUPA3k22Ncx4XXZS9o0MBH27Bo6BpNelZpS+/uh9KsNlY6bHCmJU9p8g7m3fVKn28H3KDYA5Pl/T8Z1ptDAVe0lXdQ2YoyyH2uyPIGHBZZIs2pDBS8R07+qN+E7Q=="
	msgData        = `<xml><ToUserName><![CDATA[wx5823bf96d3bd56c7]]></ToUserName>
<FromUserName><![CDATA[mycreate]]></FromUserName>
<CreateTime>1409659813</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[hello]]></Content>
<MsgId>4561255354251345929</MsgId>
<AgentID>218</AgentID>
</xml>`
)
