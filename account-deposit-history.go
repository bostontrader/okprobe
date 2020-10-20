package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
)

/*
The OKEx API has two similar endpoints that are used to get a deposit history:
1. /deposit/history
2. /deposit/history/somecurrency

The former (this probe) retrieves all deposit history for all currencies.  The later retrieves the same but filters by a specified currency.

These are superficially similar but there are subtle and nettlesome differences because the currency is specified (or not)
as part of the URL or as a query string.  These differences confound the http client and the formation of request signatures.

We have therefore chosen to deal with this as two separate endpoints.  Please also see account-deposit-history-by-cur.

*/
func ProbeAccountDepositHistory(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/deposit/history"

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

	// 2. After we've tried all the errors, it's time to build and submit the final correct request.

	// 2.1 Build a request
	req, err := standardGETReq(credentials, endPoint, queryString, baseURL)
	if err != nil {
		fmt.Println("Error building the request 2.1 : ", err)
		return
	}

	// 2.2 We expect a 2xx response
	body := TestitAPICore(httpClient, req, 200)

	// 2.3 Ensure that the prior response is parsable.
	depositHistories := make([]utils.DepositHistory, 0)
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositHistories)
	if err != nil {
		fmt.Println("Error parsing string into json 2.3: ", err)
		return
	}
}
