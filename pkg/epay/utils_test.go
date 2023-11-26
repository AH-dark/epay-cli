package epay

import (
	"github.com/AH-dark/bytestring"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateEPaySign(t *testing.T) {
	asserts := assert.New(t)

	testData := map[string]string{
		"pid":          "1001",
		"type":         "alipay",
		"out_trade_no": "20160806151343349",
		"notify_url":   "https://www.pay.com/notify_url.php",
		"return_url":   "https://www.pay.com/return_url.php",
		"name":         "VIP会员",
		"money":        "1.00",
		"clientip":     "192.168.1.100",
		"device":       "pc",
		"param":        "",
		"sign":         "4a0d8e5a6499e5de11878b87fa9ac8ba",
		"sign_type":    "md5",
	}

	asserts.Equal(testData["sign"], bytestring.BytesToString(CalculateEPaySign(testData, "123456")))
}
