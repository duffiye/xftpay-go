package qr_pre_pay

import (
	"fmt"
	"testing"
	"time"

	"github.com/duffiye/xftpay-go/xftpay"
)

func TestQrPrePay(t *testing.T) {
	qrPrePayRequest := &QrPrePayRequest{}

	qrPrePayRequest.AppID = "1052864091804139520"
	qrPrePayRequest.MerchantCode = "1130300000261"
	qrPrePayRequest.OutTradeNo = fmt.Sprintf("%v", time.Now().Nanosecond())
	qrPrePayRequest.SignType = "MD5"
	qrPrePayRequest.Amount = "1"
	qrPrePayRequest.StoreCode = "3430100000996"

	xftpay.Key = "62C433D13B689A91FAE54C5E1CC4C954"

	qrPrePayResponse, err := Pay(qrPrePayRequest)
	if err != nil {
		fmt.Printf("QrPrePay Error is %s\n", err)
	}
	fmt.Printf("QrPrePay Response is %v\n", qrPrePayResponse)

}
