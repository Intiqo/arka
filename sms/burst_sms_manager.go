package sms

import (
	"fmt"
	"github.com/adwitiyaio/arka/secrets"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	BurstSmsApikey    = "BURST_SMS_API_KEY"
	BurstSmsSecretKey = "BURST_SMS_API_SECRET_KEY"
)

type BurstSmsManager struct {
	sm        secrets.Manager
	client    *resty.Client
	apiKey    string
	apiSecret string
}

func (bs *BurstSmsManager) initialize() {
	bs.client = resty.New()
	bs.apiKey = bs.sm.GetValueForKey(BurstSmsApikey)
	bs.apiSecret = bs.sm.GetValueForKey(BurstSmsSecretKey)
}

func (bs *BurstSmsManager) SendSms(options Options) (interface{}, error) {
	requestUrl := "https://api.transmitsms.com/send-sms.json"
	method := "POST"

	recipients := strings.Join(options.Recipients, ",")

	formData := url.Values{
		"message": {options.Message},
		"to":      {recipients},
	}
	req, err := http.NewRequest(method, requestUrl, strings.NewReader(formData.Encode()))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	req.SetBasicAuth(bs.apiKey, bs.apiSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if os.Getenv("CI") != "true" {
		return bs.dispatch(req)
	}
	return nil, nil
}

func (bs *BurstSmsManager) dispatch(req *http.Request) (interface{}, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	return nil, err
}
