package qr_pre_pay

import (
	"fmt"
	"log"
	"time"

	"github.com/duffiye/xftpay-go/xftpay"
)

type Client struct {
	B   xftpay.Backend
	Key string
}

func getC() Client {
	return Client{
		xftpay.GetBackend(xftpay.APIBackend),
		xftpay.Key,
	}
}

var path string = "/transaction/unify/charge/qrPrePay"

type QrPrePayRequest struct {
	AppID        string `json:"app_id"`
	MerchantCode string `json:"merchant_code"`
	StoreCode    string `json:"store_code"`
	OperatorID   string `json:"operator_id"`
	OutTradeNo   string `json:"out_trade_no"`
	Amount       string `json:"amount"`
	Extra        string `json:"extra"`
	// Extra        struct {
	// 	QrCode2Img string `json:"qr_code_to_img"`
	// } `json:"extra"`
	SignType string `json:"sign_type"`
	Sign     string `json:"sign"`
}

type QrPrePayResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
	Data    struct {
		OutTradeNo string `json:"out_trade_no"`
		Credential string `json:"credential"`
		// Credential struct {
		// 	QrCode       string `json:"qr_code"`
		// 	QrCodeImgURL string `json:"qr_code_img_url"`
		// } `json:"credential"`
		FailureCode string `json:"failure_code"`
		FailureMsg  string `json:"failure_msg"`
		SignType    string `json:"sign_type"`
		Sign        string `json:"sign"`
	}
}

func (c Client) Pay(r *QrPrePayRequest) (s *QrPrePayResponse, err error) {
	start := time.Now()
	req, err := xftpay.GenMD5Sign(r, c.Key)
	fmt.Printf("PayRequest Sign is %s\n", req)
	if err != nil {
		if xftpay.LogLevel > 0 {
			log.Printf("Get MD5 Sign Error is : %q\n", err)
		}
	}
	paramerString, err := xftpay.JsonEncode(&req)
	fmt.Printf("JSON is %s\n", paramerString)
	if err != nil {
		if xftpay.LogLevel > 0 {
			log.Printf("PayRequest JSON Marshall Error is : %q\n", err)
		}
	}
	if xftpay.LogLevel > 2 {
		log.Printf("params of charge request to xftpay is :\n %v\n ", string(paramerString))
	}
	s = &QrPrePayResponse{}
	err = c.B.Call("POST", path, c.Key, nil, paramerString, s)
	if err != nil {
		return nil, err
	}
	if xftpay.LogLevel > 2 {
		log.Println("Pay completed in ", time.Since(start))
	}
	return s, err
}

func Pay(r *QrPrePayRequest) (s *QrPrePayResponse, err error) {
	return getC().Pay(r)
}
