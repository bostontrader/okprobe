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

func ProbeCurrencies(urlBase string, keyFile string) {

	endpoint := "/api/account/v3/currencies"
	url := urlBase + endpoint
	client := GetClient(urlBase)

	TestitStd(client, url)

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

	// 7.
	req, err := http.NewRequest("GET", url, nil)
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

	// The final correct request must have a valid signature, so build one now.
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endpoint
	encoded, _ := utils.HmacSha256Base64Signer(prehash, obj.SecretKey)

	req, err = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
		"Vary":                      "",
	}
	body := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	currenciesEntries := make([]utils.CurrenciesEntry, 0)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&currenciesEntries)
	if err != nil {
		panic(err)
	}
	fmt.Println(&currenciesEntries)
	fmt.Println(reflect.TypeOf(currenciesEntries))
}
