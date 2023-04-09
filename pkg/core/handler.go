package core

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/dthung1602/arc/pkg/resp"
)

// ---------------------------------
//	Interfaces
// ---------------------------------

type CommandHandler interface {
	Handle(req resp.Array) (resp.Resp, error)
}

type CommandHandlerFactory func(req resp.Array) (CommandHandler, error)

// ---------------------------------
//	Implementation
// ---------------------------------

func CommandHandlerFactoryImpl(arr resp.Array) (CommandHandler, error) {
	if len(arr) == 0 {
		return nil, errors.New("empty request")
	}

	cmdBS := arr[0].(resp.RespString)
	if cmdBS == nil {
		return nil, errors.New("command must be a string")
	}
	cmd := cmdBS.ToRawStr()

	handler, validCmd := cmdMapping[strings.ToUpper(cmd)]
	if !validCmd {
		return nil, fmt.Errorf("unknown command %s", cmd)
	}
	return handler, nil
}

var cmdMapping = map[string]CommandHandler{
	"COMMAND": &CommandCommandHandler{},
	"GET":     &GetCommandHandler{},
	"SET":     &SetCommandHandler{},
	"KEYS":    &KeysCommandHandler{},
}

// ---------------------------------
//	COMMAND
// ---------------------------------

type CommandCommandHandler struct{}

func (handler CommandCommandHandler) Handle(req resp.Array) (resp.Resp, error) {
	return resp.EmptyArray, nil
}

// ---------------------------------
//	GET
// ---------------------------------

type GetCommandHandler struct{}

func (handler GetCommandHandler) Handle(req resp.Array) (resp.Resp, error) {
	if len(req) != 2 {
		return nil, errors.New("wrong number of parameter for GET")
	}

	key := resp.ToByteSlice(req[1])
	if key == nil {
		return nil, errors.New("key must be of type string")
	}

	val := hashMapInstance.Get(key)
	if val == nil {
		return resp.NullVal, nil
	}
	return resp.ByteSliceToRespString(val), nil
}

// ---------------------------------
//	SET
// ---------------------------------

type SetCommandHandler struct{}

func (handler SetCommandHandler) Handle(req resp.Array) (resp.Resp, error) {
	if len(req) != 3 {
		return nil, errors.New("wrong number of parameter for SET")
	}

	key := resp.ToByteSlice(req[1])
	if key == nil {
		return nil, errors.New("key must be of type string")
	}

	val := resp.ToByteSlice(req[2])
	if key == nil {
		return nil, errors.New("value must be of type string")
	}

	hashMapInstance.Set(key, val)

	return resp.OKString, nil
}

// ---------------------------------
//	KEYS
// ---------------------------------

type KeysCommandHandler struct{}

func (handler KeysCommandHandler) Handle(req resp.Array) (resp.Resp, error) {
	if len(req) != 2 {
		return nil, errors.New("wrong number of parameter for KEYS")
	}

	rawPattern := resp.ToByteSlice(req[1])
	if rawPattern == nil {
		return nil, errors.New("pattern must be of type string")
	}

	pattern, err := regexp.CompilePOSIX("^" + strings.Replace(string(rawPattern), "*", ".*", -1) + "$")
	if err != nil {
		return nil, err
	}

	existingKeys := resp.Array{}

	for key := range hashMapInstance {
		match := pattern.FindStringIndex(key)
		if match != nil {
			existingKeys = append(existingKeys, resp.StrToRespString(key))
		}
	}

	return existingKeys, nil
}
