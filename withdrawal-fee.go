package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

func ProbeWithdrawalFee(urlBase string, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/withdrawal/fee"
	url := urlBase + endpoint
	client := GetClient(urlBase)

	if makeErrors {
		TestitStd(client, url)
	}

	// In order to proceed we need to get real credentials.  Read them from a file.
	var obj utils.Credentials

	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Print(err)
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	if makeErrors {
		// 7.
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

		// 8.
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

		// 9.
		req, err = http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30013(), 401) // Invalid Sign
	}

	// Requests after this point require a valid signature.

	// 10. Extraneous parameters appear to be ignored, the full list is returned, and two extra headers appear: Vary and Strict-Transport-Security.  Status = 200.  For example: ?catfood and ?catfood=yum.
	// But this is of minimal importance so don't bother trying to test for this.  We certainly don't care to mimic this behavior in the catbox.

	if makeErrors {
		// 11. Request an invalid currency
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		invalid_param := "catfood"
		params := "?currency=" + invalid_param
		prehash := timestamp + "GET" + endpoint + params
		encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
		req, _ := http.NewRequest("GET", url+params, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30031(invalid_param), 400)
	}

	// 12. Request a single valid currency
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	valid_param := "BTC"
	params := "?currency=" + valid_param
	prehash := timestamp + "GET" + endpoint + params
	encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
	req, err := http.NewRequest("GET", url+params, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	withdrawlFees := make([]utils.WithdrawalFee, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&withdrawlFees)
	if err != nil {
		panic(err)
	}
	fmt.Println(&withdrawlFees)
	fmt.Println(reflect.TypeOf(withdrawlFees))

	// 13. Don't request any currency, so by default get them all
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "GET" + endpoint
	encoded, _ = utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
		"Vary":                      "",
	}
	body = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	withdrawlFees = make([]utils.WithdrawalFee, 0)
	dec = json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&withdrawlFees)
	if err != nil {
		panic(err)
	}
	fmt.Println(&withdrawlFees)
	fmt.Println(reflect.TypeOf(withdrawlFees))
}

func catMap(a, b map[string]string) map[string]string {
	var n = map[string]string{}
	for k, v := range a {
		n[k] = v
	}
	for k, v := range b {
		n[k] = v
	}

	return n
}
