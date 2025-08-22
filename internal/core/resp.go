package core

import (
	"bytes"
	"duchm1606/rogo/internal/constant"
	"errors"
	"fmt"
	"strings"
)

const CRLF string = "\r\n"

func Encode(value interface{}, isSimpleString bool) []byte {
	switch v := value.(type){ 
	case string:
		if isSimpleString {
			return []byte(fmt.Sprintf("+%s%s", v, CRLF))
		}
		return []byte(fmt.Sprintf("$%d%s%s%s", len(v), CRLF, v, CRLF))
	case int64, int32, int16, int8, int:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	case []string:
		return encodeStringArray(value.([]string))
	case [][]string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, sa := range value.([][]string) {
			buf.Write(encodeStringArray(sa))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([][]string)), buf.Bytes()))
	case []interface{}:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, x := range value.([]interface{}) {
			buf.Write(Encode(x, false))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(value.([]interface{})), buf.Bytes()))
	default:
		return constant.RespNil
	}
}

func encodeString(s string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
}

func encodeStringArray(sa []string) []byte {
	var b []byte
	buf := bytes.NewBuffer(b)
	for _, s := range sa {
		buf.Write(encodeString(s))
	}
	return []byte(fmt.Sprintf("*%d\r\n%s", len(sa), buf.Bytes()))
}

func ParseCommand(data []byte) (*Command, error) {
	value, err := Decode(data)
	if err != nil {
		return nil, err
	}

	array := value.([]interface{})
	tokens := make([]string, len(array))
	for i := range tokens {
		tokens[i] = array[i].(string)
	}
	res := &Command{Cmd: strings.ToUpper(tokens[0]), Args: tokens[1:]}
	return res, nil
}

func Decode(data []byte) (interface{}, error) {
	res, _, err := DecodeOne(data)
	return res, err
}

// Decode the first obbject in the data
func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("empty data")
	}
	switch data[0] {
	case '*':
		return decodeArray(data)
	case '$':
		return decodeBulkString(data)
	case '+':
		return decodeSimpleString(data)
	case ':':
		return decodeInteger(data)
	case '-':
		return decodeError(data)

	}
	return nil, 0, nil
}

// $5\r\nhello\r\n => 5, 4
func readLen(data []byte) (int, int) {
	res, pos, _ := decodeInteger(data)

	return int(res), pos
}

// Examples:
//   `*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n` -> returns []interface{} with "foo" and "bar"
func decodeArray(data []byte) (interface{}, int, error) {
	length, pos := readLen(data)
	var res []interface{} = make([]interface{}, length)

	for i := range res {
		elem, delta, err := DecodeOne(data[pos:])

		if err != nil {
			return nil, 0, err
		}
		res[i] = elem
		pos += delta
	}

	return res, pos, nil
}

// $5\r\nhello\r\n => "hello"
func decodeBulkString(data []byte) (interface{}, int, error) {
	length, pos := readLen(data)
	return string(data[pos:(pos + length)]), pos + length + 2, nil
}

func decodeSimpleString(data []byte) (interface{}, int, error) {
	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

// Parses a RESP integer from the given byte slice.
//
// RESP integers are encoded as a CRLF-terminated string that represents a signed,
// base-10, 64-bit integer. The format follows the pattern: 
//	`:[<+|->]<value>\r\n`
// Examples:
//   ":123\r\n"     -> returns int64(123)
//   ":-456\r\n"    -> returns int64(-456)
//   ":+789\r\n"    -> returns int64(789)
func decodeInteger(data []byte) (int64, int, error) {
	var res int64 = 0
	var sign int64 = 1
	pos := 1
	if data[pos] == '-' {
		sign = -1
		pos++
	}
	if data[pos] == '+' {
		sign = 1
		pos++
	}
	for data[pos] != CRLF[0] {
		res = res*10 + int64(data[pos] - '0')
		pos++
	}
	return res * sign, pos + 2, nil
}

func decodeError(data []byte) (interface{}, int, error) {
	return decodeSimpleString(data)
}