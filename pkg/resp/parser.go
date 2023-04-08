package resp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

type Parser struct{}

func (p Parser) Parse(r *bufio.Reader) (Resp, error) {
	ch, err := r.Peek(1)
	if err != nil {
		return nil, err
	}

	switch ch[0] {
	case '$':
		return p.parseBlobString(r)
	case '+':
		return p.parseSimpleString(r)
	case '-':
		return p.parseSimpleError(r)
	case ':':
		return p.parseNumber(r)
	case '*':
		return p.parseArray(r)
	default:
		return nil, fmt.Errorf("unreconized type %b", ch[0])
	}
}

// TODO handle when payload is too large
func readTillCRNL(r *bufio.Reader) ([]byte, error) {
	var data [][]byte
	for {
		buff, err := r.ReadBytes('\r')
		if err != nil {
			return nil, err
		}
		nextByte, err := r.Peek(1)
		if err != nil {
			return nil, err
		}
		if nextByte[0] == '\n' {
			data = append(data, buff[:len(buff)-1])
			_, _ = r.Discard(1)
			break
		}
		data = append(data, buff)
	}
	return bytes.Join(data, nil), nil
}

func readLen(r *bufio.Reader) (int, error) {
	lenStr, err := readTillCRNL(r)
	if err != nil {
		return -1, err
	}
	n, err := strconv.Atoi(string(lenStr))
	if err != nil {
		return -1, err
	}
	if n < 0 {
		return -1, errors.New("blob string length cannot be negative")
	}
	return n, nil
}

func (p Parser) parseBlobString(r *bufio.Reader) (BlobString, error) {
	_, _ = r.ReadByte()

	blobStrLen, err := readLen(r)
	if err != nil {
		return nil, err
	}

	blobString := make(BlobString, blobStrLen)
	n, err := r.Read(blobString)
	if err != nil {
		return nil, err
	}
	if n != blobStrLen {
		return nil, errors.New("missing bytes when reading blob string")
	}

	nextBytes, err := r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if len(nextBytes) != 2 || nextBytes[0] != '\r' {
		return nil, errors.New("missing crnl at the end of blob string")
	}

	return blobString, nil
}

func (p Parser) parseSimpleString(r *bufio.Reader) (SimpleString, error) {
	return nil, nil
}

func (p Parser) parseSimpleError(r *bufio.Reader) (SimpleError, error) {
	return nil, nil
}

func (p Parser) parseNumber(r *bufio.Reader) (Number, error) {
	return 0, nil
}

func (p Parser) parseArray(r *bufio.Reader) (Array, error) {
	_, _ = r.ReadByte()

	arrayLen, err := readLen(r)
	if err != nil {
		return nil, err
	}

	arr := make(Array, arrayLen)

	for i := 0; i < arrayLen; i++ {
		ele, err := p.Parse(r)
		if err != nil {
			return nil, err
		}
		arr[i] = ele
	}

	return arr, nil
}
