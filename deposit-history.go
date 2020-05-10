package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"reflect"
	"time"
)

/*
The OKEx API has two similar endpoint that are used to get a deposit history:
1. /deposit/history
2. /deposit/history/somecurrency

I would like to generalize this to /deposit/history/:currencyid (optional) but apparently the default go router doesn't support this.  Even if I could easily pick out the optional :currencyid the API still returns all deposit history items, not filtered by :currencyid.

This is not a priority issue so we'll only focus on using the /deposit/history variation at this time.

More particularly:

A. /deposit/history (no :currencyid) works just fine, but it gets the history of all deposits, not filtered to a particular currency.

B. If :currencyid == an invalid currency then we get Err30031 as we might expect.

C. If :currencyid == a valid currency then the call still returns all deposit records, as if we had not specified any :currencyid at all.
*/
func ProbeDepositHistory(urlBase string, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/deposit/history"
	url := urlBase + endpoint
	client := GetClient(urlBase)
	obj := getCredentials(keyFile)

	if makeErrors {
		TestitStd(client, url)
	}

	if makeErrors {
		// 7.
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

		// 8.
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

		// 9.
		req, _ = http.NewRequest("GET", url, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", "wrong")
		req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
		req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30013(), 401) // Invalid Sign
	}

	// Requests after this point require a valid signature.

	// Don't bother
	/*if makeErrors {
		// 11. Request an invalid currency
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		invalid_currency := "catfood"
		params := "/" + invalid_currency
		prehash := timestamp + "GET" + endpoint + params
		encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
		req, _ := http.NewRequest("GET", url+params, nil)
		req.Header.Add("OK-ACCESS-KEY", obj.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
		Testit4xx(client, req, utils.ExpectedResponseHeaders, utils.Err30031(invalid_currency), 400)
	}*/

	// Don't bother
	// 12. Request a single valid currency.  This actually retrieves all records, not just limited to the specified currency.
	/*timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	valid_currency := "btc"
	params := "?currency=" + valid_currency
	prehash := timestamp + "GET" + endpoint + params
	encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
	req, err := http.NewRequest("GET", url+params, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
		"Vary": "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	depositHistories := make([]utils.DepositHistory, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositHistories)
	if err != nil {
		panic(err)
	}
	fmt.Println(&depositHistories)
	fmt.Println(reflect.TypeOf(depositHistories))*/

	// 13. Don't specify any currency, so by default get history records for all
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endpoint
	encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)

	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
		"Vary":                      "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	depositHistories := make([]utils.DepositHistory, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositHistories)
	if err != nil {
		panic(err)
	}
	fmt.Println(&depositHistories)
	fmt.Println(reflect.TypeOf(depositHistories))
}
