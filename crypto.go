package workwx

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sort"
	"strings"
)

type Crypto struct {
	token string
	iv    []byte
	block cipher.Block
	id    []byte //cropid or suiteid
}

func paddingPKCS7(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func unpaddingPKCS7(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])

	return origData[:(length - unpadding)]
}

func NewCrypto(token, id string, key []byte) (*Crypto, error) {
	if block, err := aes.NewCipher(key); err != nil {
		return nil, err
	} else {
		return &Crypto{token, key[:16], block, []byte(id)}, nil
	}
}

func (p *Crypto) Encrypt(message []byte) string {
	messageSize := len(message)
	count := 20 + messageSize + len(p.id)
	buf := make([]byte, count)
	rand.Read(buf[:16])

	binary.BigEndian.PutUint32(buf[16:20], uint32(messageSize))
	copy(buf[20:20+messageSize], message)
	copy(buf[20+messageSize:], p.id)
	buf = paddingPKCS7(buf, aes256BlockSize)

	data := make([]byte, len(buf))
	aes := cipher.NewCBCEncrypter(p.block, p.iv)
	aes.CryptBlocks(data, buf)
	base64 := base64.StdEncoding

	return base64.EncodeToString(data)
}

func (p *Crypto) Decrypt(message string) ([]byte, error) {
	base64 := base64.StdEncoding
	data, err := base64.DecodeString(message)

	if err != nil {
		return nil, err
	} else if len(data) < 16+4+len(p.id) {
		return nil, ErrMessage
	} else {
		buf := make([]byte, len(data))
		dec := cipher.NewCBCDecrypter(p.block, p.iv)
		dec.CryptBlocks(buf, data)
		buf = unpaddingPKCS7(buf)
		msgLen := int(binary.BigEndian.Uint32(buf[16:20]))
		idPos := 16 + 4 + msgLen

		if len(buf) != idPos+len(p.id) || bytes.Compare(p.id, buf[idPos:]) != 0 {
			return nil, ErrMessage
		} else {
			return buf[20:idPos], nil
		}
	}
}

func (p *Crypto) Signature(timeStamp, nonce, message string) string {
	strs := []string{
		p.token, timeStamp, nonce, message,
	}
	sort.Strings(strs)
	str := strings.Join(strs, "")
	data := sha1.Sum([]byte(str))

	return hex.EncodeToString(data[:])
}

const (
	aes256BlockSize = 32
)

var (
	ErrMessage = errors.New("message error")
)
