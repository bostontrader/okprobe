package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
)

func ProbeAccountWallet(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/wallet"

	// 1.1 Read and parse credentials file
	credentials := getCredentials(credentialsFile)

	// 1.2 Obtain an http client
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)

	// 1.3 If we want to test header/credentials errors.
	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	// 2. After we've tried all the errors, it's time to build and submit the final correct request.

	// 2.1 Build a request
	req := buildGETRequest(credentials, endPoint, "", baseURL)

	// 2.2 We expect a 2xx response
	body := TestitAPICore(httpClient, req, 200)

	// 2.3 Ensure that the prior response is parsable.
	walletEntries := make([]utils.WalletEntry, 0)
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	err := dec.Decode(&walletEntries)
	if err != nil {
		fmt.Println("Error parsing string into json 2.3: ", err)
		return
	}
}
