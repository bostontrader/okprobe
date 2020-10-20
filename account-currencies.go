package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
)

func ProbeAccountCurrencies(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/currencies"

	// 1.1 Read and parse credentials file
	credentials, err := getCredentials(credentialsFile)
	if err != nil {
		fmt.Println("Error obtaining the credentials 1.1: ", err)
		return
	}

	// 1.2 Obtain an http client
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)

	// 1.3 If we want to test header/credentials errors.
	if makeErrorsCredentials {
		TestitCredentialsErrors(httpClient, url, credentials)
	}

	if forReal {
		// 2.1 Build a request
		req, err := standardGETReq(credentials, endPoint, "", baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.1 : ", err)
			return
		}

		// 2.2 We expect a 2xx response
		body := TestitAPICore(httpClient, req, 200)

		// 2.3 Ensure that the prior response is parsable.
		currencyEntries := make([]utils.CurrenciesEntry, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err = dec.Decode(&currencyEntries)
		if err != nil {
			fmt.Println("Error parsing string into json 2.3: ", err)
			return
		}
	}
}
