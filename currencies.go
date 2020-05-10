package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"reflect"
	"time"
)

func ProbeCurrencies(urlBase string, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/currencies"
	url := urlBase + endpoint
	client := GetClient(urlBase)
	credentials := getCredentials(keyFile)

	if makeErrors {
		TestitStd(client, url, credentials, utils.ExpectedResponseHeaders)
	}

	// The final correct request must have a valid signature, so build one now.
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endpoint
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
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
