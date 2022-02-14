package ovopark

import (
	"context"
	"testing"
)

type SWPFGetAllHotspotShopsReq struct {
	CommonReqFields
	CorpId   int64  `json:"orgid" url:"orgid"`
	DepId    int    `json:"depId" url:"depId"`
	Page     int    `json:"pageNumber" url:"pageNumber"`
	Limit    int    `json:"pageSize" url:"pageSize"`
	ShopId   string `json:"shopId" url:"shopId"`
	SWCorpId int64  `json:"ovoparkEnterpriseId" url:"ovoparkEnterpriseId"`
}

type SWPFGetAllHotspotShopsResp struct {
	CommonRespFields
	Data struct {
		Pages       int `json:"total"`
		Total       int `json:"all"`
		CurrentPage int `json:"pageNumber"`
		Limit       int `json:"pageSize"`
		List        struct {
			Devices []struct {
				Online     int    `json:"online"`
				DeviceUrl  string `json:"deviceUrl"`
				DeviceName string `json:"deviceName"`
				Mac        string `json:"mac"`
			} `json:"devices"`
			DepId   int    `json:"depId"`
			DepName string `json:"depName"`
			Nums    int    `json:"nums"`
		}
	} `json:"data"`
}

func TestAPIRequest(t *testing.T) {
	client := NewOVOPark(0, "fake-appid", "fake-akeyid", "fake-akeysec")
	client.SetAuthenticator("fake-token")
	req := SWPFGetAllHotspotShopsReq{CorpId: 0}
	resp := SWPFGetAllHotspotShopsResp{}
	err := client.APIRequest(context.Background(), "open.passengerflow.getAllHotspotShops", &req, &resp)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
