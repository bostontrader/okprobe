package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"reflect"
	"time"
)

func ProbeDepositAddress(urlBase string, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/deposit/address"
	url := urlBase + endpoint
	client := GetClient(urlBase)
	credentials := getCredentials(keyFile)

	if makeErrors {
		TestitStd(client, url, credentials, utils.ExpectedResponseHeaders)
	}

	// Requests after this point require a valid signature.

	// 10. Extraneous parameters appear to be ignored, the full list is returned, and two extra headers appear: Vary and Strict-Transport-Security.  Status = 200.  For example: ?catfood and ?catfood=yum.
	// But this is of minimal importance so don't bother trying to test for this.  We certainly don't care to mimic this behavior in the catbox.

	if makeErrors {
		// 11. Don't request any currency
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash := timestamp + "GET" + endpoint
		encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30023(), 400)

		// 12. Request an invalid currency
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		invalid_param := "catfood"
		params := "?currency=" + invalid_param
		prehash = timestamp + "GET" + endpoint + params
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
		req, _ = http.NewRequest("GET", url+params, nil)
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30031(invalid_param), 400)
	}

	// 13. Request a valid currency
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	valid_param := "BTC"
	params := "?currency=" + valid_param
	prehash := timestamp + "GET" + endpoint + params
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("GET", url+params, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	depositAddress := make([]utils.DepositAddress, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositAddress)
	if err != nil {
		panic(err)
	}
	fmt.Println(&depositAddress)
	fmt.Println(reflect.TypeOf(depositAddress))
}
