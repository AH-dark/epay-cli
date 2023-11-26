package epay

import (
	"crypto/subtle"
	"fmt"
	"sort"
	"strings"

	"github.com/AH-dark/bytestring"
	"github.com/samber/lo"

	"github.com/AH-dark/epay-cli/pkg/utils"
)

// GenerateParams 生成加签参数
func GenerateParams(params map[string]string, secret string) map[string]string {
	params["sign"] = bytestring.BytesToString(CalculateEPaySign(params, secret))
	params["sign_type"] = "MD5"
	return params
}

func CalculateEPaySign(mapData map[string]string, secret string) []byte {
	// sort keys
	keys := lo.Keys(mapData)
	sort.Strings(keys)

	combinedData := ""
	for _, k := range keys {
		if k == "sign" || k == "sign_type" || lo.IsEmpty(mapData[k]) {
			continue
		}

		combinedData += fmt.Sprintf("%s=%s&", k, mapData[k])
	}

	combinedData = strings.TrimSuffix(combinedData, "&")

	return utils.MD5(combinedData + secret)
}

func CheckEPaySign(mapData map[string]string, secret string, sign []byte) bool {
	return subtle.ConstantTimeCompare(CalculateEPaySign(mapData, secret), sign) == 1
}
