package charge

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

var path string = "/transaction/unify/charge/pay"

type PayRequest struct {
	AppID        string `json:"app_id"`
	MerchantCode string `json:"merchant_code"`
	StoreCode    string `json:"store_code"`
	LimitPay     string `json:"limit_pay"`
	OutTradeNo   string `json:"out_trade_no"`
	Channel      string `json:"channle"`
	Product      string `json:"product"`
	ClientIP     string `json:"client_ip"`
	Amount       string `json:"amount"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	Descryption  string `json:"description"`
	Extra        string `json:"extra"`
	NotifyURL    string `json:"notify_url"`
	TimeStart    string `json:"time_start"`
	TimeExpire   string `json:"time_expire"`
	SignType     string `json:"sign_type"`
	Sign         string `json:"sign"`
}

type PayResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
	Data    struct {
		ID           string `json:"id"`
		OutTradeNo   string `json:"out_trade_no"`
		State        string `json:"state"`
		Credential   string `json:"credential"`
		ThirdTradeNo string `json:"third_trade_no"`
		FailureCode  string `json:"failure_code"`
		FailureMsg   string `json:"failure_msg"`
		SignType     string `json:"sign_type"`
		Sgin         string `json:"sign"`
	} `json:"data"`
}

func (c Client) Pay(r *PayRequest) (s *PayResponse, err error) {
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
	s = &PayResponse{}
	err = c.B.Call("POST", path, c.Key, nil, paramerString, s)
	if err != nil {
		return nil, err
	}
	if xftpay.LogLevel > 2 {
		log.Println("Pay completed in ", time.Since(start))
	}
	return s, err
}

func Pay(r *PayRequest) (s *PayResponse, err error) {
	return getC().Pay(r)
}
