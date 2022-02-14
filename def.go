package ovopark

type CommonReqFields struct {
	APIName      string `json:"_mt" url:"_mt"`
	APIVersion   string `json:"_version" url:"_version"`
	ReqMethod    string `json:"_requestMode" url:"_requestMode"`
	PlatformID   string `json:"_aid" url:"_aid"`
	AccessKeyId  string `json:"_akey" url:"_akey"`
	CryptoMethod string `json:"_sm" url:"_sm"`
	RequestTime  string `json:"_timestamp" url:"_timestamp"`
	ReqToken     string `json:"_sig" url:""`
	RespFormat   string `json:"_format" url:"_format"`
}

type CommonRespFields struct {
	Status struct {
		ReqId        string `json:"cid"`
		ErrCode      int    `json:"code"`
		ErrCodeStr   string `json:"codename"`
		ServerTime   int64  `json:"systime"`
		GatewayParam struct {
			IsGatewayReturn bool `json:"isGatewayReturn"`
		} `json:"gatewayParam"`
	} `json:"stat"`
	ErrMsg string `json:"result"`
}

type CommonResp struct {
	CommonRespFields
	Data string `json:"data,omitempty"`
}
