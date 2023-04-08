package resp

import (
	"bytes"
	"strconv"
	"strings"
)

// ---------------------------------
//	Interfaces
// ---------------------------------

type Resp interface {
	Resp() []byte
	String() string
}

// ---------------------------------
//	Util functions
// ---------------------------------

func simpleResp(s []byte, marker byte) []byte {
	totalLen := 3 + len(s)
	result := make([]byte, totalLen)

	result[0] = marker
	copy(result[1:], s)
	result[totalLen-2] = '\r'
	result[totalLen-1] = '\n'
	return result
}

// ---------------------------------
//	BlobString $
// ---------------------------------

type BlobString []byte

func (s BlobString) Resp() []byte {
	lengthStr := []byte(strconv.Itoa(len(s)))
	lls := len(lengthStr)
	totalLen := 5 + lls + len(s)
	result := make([]byte, totalLen)

	result[0] = '$'
	copy(result[1:], lengthStr)
	result[lls+1] = '\r'
	result[lls+2] = '\n'
	copy(result[lls+3:], s)
	result[totalLen-2] = '\r'
	result[totalLen-1] = '\n'
	return result
}

func (s BlobString) String() string {
	return "BlobString(" + string([]byte(s)) + ")"
}

// ---------------------------------
//	SimpleString +
// ---------------------------------

type SimpleString []byte

func (s SimpleString) Resp() []byte {
	return simpleResp(s, '+')
}

func (s SimpleString) String() string {
	return "SimpleString"
}

// ---------------------------------
//	SimpleError -
// ---------------------------------

type SimpleError []byte

func (e SimpleError) Resp() []byte {
	return simpleResp(e, '-')
}

func (e SimpleError) String() string {
	return "SimpleError"
}

// ---------------------------------
//	Number :
// ---------------------------------

type Number int64

func (n Number) Resp() []byte {
	return simpleResp([]byte(strconv.FormatInt(int64(n), 10)), ':')
}

func (n Number) String() string {
	return "Number"
}

// ---------------------------------
//	Array *
// ---------------------------------

type Array []Resp

func (a Array) Resp() []byte {
	resps := make([][]byte, len(a)+3)

	resps[0] = []byte{'*'}
	resps[1] = []byte(strconv.Itoa(len(a)))
	resps[2] = []byte{'\r', '\n'}

	for i, item := range a {
		resps[i+3] = item.Resp()
	}

	return bytes.Join(resps, nil)
}

func (a Array) String() string {
	var eleStrs []string

	for i := 0; i < len(a); i++ {
		blobStr, isBlobStr := a[i].(BlobString)
		if isBlobStr {
			eleStrs = append(eleStrs, blobStr.String())
		}
	}
	return "Array[\n    " + strings.Join(eleStrs, ",\n    ") + "\n]"
}
