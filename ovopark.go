package ovopark

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"reflect"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/google/go-querystring/query"
)

type OVOPark struct {
	domain        string
	platformID    string
	accessKeyId   string
	accessKeySec  string
	authenticator string
	corpId        int
}

func NewOVOPark(corpId int, aid, aKeyId, aKeySec string) *OVOPark {
	s := new(OVOPark)
	s.platformID = aid
	s.accessKeyId = aKeyId
	s.accessKeySec = aKeySec
	s.corpId = corpId
	s.domain = "https://cloudapi.ovopark.com/cloud.api"
	return s
}

func (s *OVOPark) SetDomain(domain string) {
	s.domain = domain
}

// SetAuthenticator 设置auth token
func (s *OVOPark) SetAuthenticator(token string) {
	s.authenticator = token
}

func (s *OVOPark) getAuthToken() map[string]string {
	var header = make(map[string]string, 0)
	if len(s.authenticator) > 0 {
		header["authenticator"] = s.authenticator
	}
	return header
}

func (s *OVOPark) APIRequest(ctx context.Context, apiName string, req, response interface{}, optional ...string) (err error) {
	_ = checkRequest(req)
	ReqValue := reflect.ValueOf(req)
	ReqStructValue := ReqValue.Elem()
	var (
		ver    = "v1"
		method = "POST"
	)
	ReqStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex).Set(reflect.ValueOf(apiName))
	if len(optional) > 0 {
		ver = optional[0]
	}
	if len(optional) > 1 {
		method = optional[1]
	}
	ReqStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 1).Set(reflect.ValueOf(ver))
	ReqStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 2).Set(reflect.ValueOf(method))
	ErrorStructValue, ErrorErrCodeValue := checkResponse(response)
	v, _ := query.Values(req)
	formStr := s.buildCommonParams(v)
	var (
		bodyByte  []byte
		jsonqNode *jsonquery.Node
	)
	_, bodyByte, err = doHttpFormReq(ctx, s.domain, formStr, "POST", s.getAuthToken())
	if err != nil {
		return
	}
	jsonqNode, err = jsonquery.Parse(bytes.NewReader(bodyByte))
	if err != nil {
		return
	}
	dataNode := jsonquery.FindOne(jsonqNode, "data")
	if len(dataNode.InnerText()) == 0 {
		var resp2 CommonResp
		err = decodeJSONHttpResponse(bytes.NewReader(bodyByte), &resp2)
		if err != nil {
			return
		}
		ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex).Set(reflect.ValueOf(resp2.Status.ReqId))
		ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 1).Set(reflect.ValueOf(resp2.Status.ErrCode))
		ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 2).Set(reflect.ValueOf(resp2.Status.ErrCodeStr))
		ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 3).Set(reflect.ValueOf(resp2.Status.ServerTime))
		ErrorStructValue.Field(errorErrCodeIndex).Field(errorErrCodeIndex + 4).Field(errorErrCodeIndex).Set(reflect.ValueOf(resp2.Status.GatewayParam.IsGatewayReturn))
		ErrorStructValue.Field(errorErrMsgIndex).Set(reflect.ValueOf(resp2.ErrMsg))
	} else {
		err = decodeJSONHttpResponse(bytes.NewReader(bodyByte), response)
	}
	switch errCode := ErrorErrCodeValue.Int(); errCode {
	case ErrCodeOK:
		return
	default:
		errMsg := ErrorStructValue.Field(errorErrMsgIndex).String()
		err = errors.New(errMsg)
		return
	}
}

func (s *OVOPark) buildCommonParams(req url.Values) string {
	req.Set("_aid", s.platformID)
	req.Set("_akey", s.accessKeyId)
	if len(req.Get("_sm")) == 0 {
		req.Set("_sm", "md5")
	}
	req.Set("_timestamp", msTimestamp2SecondStr())
	if len(req.Get("_format")) == 0 {
		req.Set("_format", "json")
	}
	sigTmpStr := req.Encode()
	sigTmpOutStr := strings.ReplaceAll(sigTmpStr, "&", "")
	sigTmpOutStr = strings.ReplaceAll(sigTmpOutStr, "=", "")
	sigTmpOutStr, _ = url.QueryUnescape(sigTmpOutStr)
	sig := s.buildSignature(req.Get("_sm"), sigTmpOutStr)
	req.Set("_sig", sig)
	return req.Encode()
}

func (s *OVOPark) buildSignature(sm, in string) (out string) {
	inTemp := s.accessKeySec + in + s.accessKeySec
	if sm == "md5" {
		out, _ = md5x(inTemp)
	} else {
		out, _ = sha1x(inTemp)
	}
	out = strings.ToUpper(out)
	return
}
