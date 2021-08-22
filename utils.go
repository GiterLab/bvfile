package bvfile

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/GiterLab/gomathbits"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// DecodeToUInt 解码bin为 uint8/uint16/uint32
func DecodeToUInt(b []byte, bits int) (uint64, error) {
	if len(b) == 1 && bits == BitsUint8 {
		return uint64(b[0]), nil
	}
	if len(b)%2 != 0 {
		return 0, errors.New("array is not even in length")
	}
	hexStrSrc := hex.EncodeToString(b)
	hexStrDes := ""
	for i := 0; i < len(hexStrSrc); i++ {
		if i%2 != 0 {
			hexStrDes = hexStrDes + string(hexStrSrc[i])
		}
	}
	return gomathbits.ParseUInt(reverse(hexStrDes), bits)
}

func decodeGBK(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// DecodeToString 解码为字符串
func DecodeToString(b []byte, bits int) string {
	if len(b) == 0 {
		return ""
	}
	b, err := decodeGBK(b)
	if err != nil {
		return ""
	}
	index := bytes.IndexByte(b, 0)
	b = b[0:index]
	return strings.Trim(string(b), "\u0000")
}
