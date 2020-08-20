package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"os"
	"strings"
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
		err := TestitCredentialsErrors(httpClient, url, credentials, utils.ExpectedResponseHeaders)
		if err != nil {
			os.Exit(1)
		}
	}

	// 2. After we've tried all the errors, it's time to build and submit the final correct request.

	// 2.1 Build a request
	req, err := standardGETReq(credentials, endPoint, "", baseURL)
	if err != nil {
		fmt.Println("Error building the request 2.1 : ", err)
		return
	}

	// 2.2 We expect a 2xx response
	body, err := TestitAPI2xx(httpClient, req, utils.ExpectedResponseHeaders)
	if err != nil {
		fmt.Println("Error invoking the API 2.2: ", err)
		return
	}
	fmt.Println(body)

	// 2.3 Ensure that the prior response is parsable.
	currencyEntries := make([]utils.CurrenciesEntry, 0)
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&currencyEntries)
	if err != nil {
		fmt.Println("Error parsing string into json 2.3: ", err)
		return
	}

}
