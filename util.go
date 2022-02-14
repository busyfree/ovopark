package ovopark

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"
)

const (
	ErrCodeOK = 0
)

var (
	errorType   = reflect.TypeOf(CommonRespFields{})
	commReqType = reflect.TypeOf(CommonReqFields{})
)

const (
	errorErrCodeIndex = 0
	errorErrMsgIndex  = 1
)

func decodeJSONHttpResponse(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func checkResponse(response interface{}) (ErrorStructValue, ErrorErrCodeValue reflect.Value) {
	responseValue := reflect.ValueOf(response)
	if responseValue.Kind() != reflect.Ptr {
		panic("the type of response is incorrect")
	}
	responseStructValue := responseValue.Elem()
	if responseStructValue.Kind() != reflect.Struct {
		panic("the type of response is incorrect")
	}

	if t := responseStructValue.Type(); t == errorType {
		ErrorStructValue = responseStructValue
	} else {
		if t.NumField() == 0 {
			panic("the type of response is incorrect")
		}
		v := responseStructValue.Field(0)
		if v.Type() != errorType {
			panic("the type of response is incorrect")
		}
		ErrorStructValue = v
	}
	ErrorErrCodeValue = ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrMsgIndex)
	return
}

func checkRequest(request interface{}) (CommReqStructValue reflect.Value) {
	reqValue := reflect.ValueOf(request)
	if reqValue.Kind() != reflect.Ptr {
		panic("the type of request is incorrect")
	}
	reqStructValue := reqValue.Elem()
	if reqStructValue.Kind() != reflect.Struct {
		panic("the type of request is incorrect")
	}
	if t := reqStructValue.Type(); t == commReqType {
		CommReqStructValue = reqStructValue
	} else {
		if t.NumField() == 0 {
			panic("the type of request is incorrect")
		}
		v := reqStructValue.Field(0)
		if v.Type() != commReqType {
			panic("the type of request is incorrect")
		}
		CommReqStructValue = v
	}
	return
}

func msTimestamp2SecondStr() string {
	return time.Now().Format("20060102150405")
}

func md5x(data string) (string, error) {
	h := md5.New()
	if _, err := h.Write([]byte(data)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func sha1x(data string) (string, error) {
	h := sha1.New()
	_, err := io.Copy(h, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
