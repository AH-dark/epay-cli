package mapi

import (
	"fmt"
	"strconv"

	"github.com/urfave/cli/v2"

	"github.com/AH-dark/bytestring"
	"github.com/AH-dark/epay-cli/actions/factory"
	"github.com/AH-dark/epay-cli/pkg/epay"
	"github.com/AH-dark/epay-cli/pkg/utils"
)

type service struct {
}

func NewService() factory.ActionService {
	return &service{}
}

func (svc *service) getSign(c *cli.Context) []byte {
	sign := epay.CalculateEPaySign(map[string]string{
		"pid":          strconv.Itoa(c.Int("pid")),
		"type":         c.String("type"),
		"out_trade_no": c.String("trade-no"),
		"notify_url":   c.String("notify-url"),
		"return_url":   c.String("return-url"),
		"name":         c.String("name"),
		"money":        c.String("money"),
		"clientip":     c.String("client-ip"),
		"device":       c.String("device"),
		"param":        c.String("param"),
	}, c.String("secret"))

	return sign
}

func (svc *service) Do(c *cli.Context) error {
	client, err := epay.NewClient(&epay.Config{
		PartnerID: c.Int("pid"),
		AppSecret: c.String("secret"),
		Endpoint:  c.String("endpoint"),
	})
	if err != nil {
		return err
	}

	sign := svc.getSign(c)
	fmt.Println("Sign:", sign)

	url, args, err := client.MApiSubmit(c.Context, &epay.MApiSubmitArgs{
		Type:       epay.PaymentType(c.String("type")),
		OutTradeNo: c.String("trade-no"),
		Name:       c.String("name"),
		Money:      c.String("money"),
		NotifyUrl:  c.String("notify-url"),
		ReturnUrl:  utils.EmptyPtr(c.String("return-url")),
		ClientIP:   c.String("client-ip"),
		Device:     utils.EmptyPtr(epay.DeviceType(c.String("device"))),
		Param:      utils.EmptyPtr(c.String("param")),
		Sign:       bytestring.BytesToString(sign),
		SignType:   "MD5",
	})
	if err != nil {
		return err
	}

	fmt.Println("URL:", url)
	fmt.Println("Args:", args)

	return nil
}
