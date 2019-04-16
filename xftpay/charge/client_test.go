package charge

import (
	"fmt"
	"testing"
	"time"

	"github.com/duffiye/xftpay-go/xftpay"
)

func TestPay(t *testing.T) {
	p := &PayRequest{}

	p.AppID = "1052864091804139520"
	p.MerchantCode = "1130300000261"
	p.OutTradeNo = fmt.Sprintf("%v", time.Now().Nanosecond())
	p.SignType = "MD5"
	p.Amount = "1"
	p.StoreCode = "3430100000996"

	xftpay.Key = "62C433D13B689A91FAE54C5E1CC4C954"

	payResponse, err := Pay(p)
	if err != nil {
		fmt.Printf("QrPrePay Error is %s\n", err)
	}
	fmt.Printf("QrPrePay Response is %v\n", payResponse)
}
