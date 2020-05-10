package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"reflect"
	"time"
)

func ProbeWallet(urlBase, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/wallet"
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

	// The final correct request must have a valid signature, so build one now.
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
	}

	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	walletEntries := make([]utils.WalletEntry, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&walletEntries)
	if err != nil {
		panic(err)
	}
	fmt.Println(&walletEntries)
	fmt.Println(reflect.TypeOf(walletEntries))
}
