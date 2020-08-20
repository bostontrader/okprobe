package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"os"
	"strings"
)

/*
The OKEx API has two similar endpoints that are used to get a deposit history:
1. /deposit/history
2. /deposit/history/somecurrency

The former retrieves all deposit history for all currencies.  The later (this probe) retrieves the same but filters by a specified currency.

These are superficially similar but there are subtle and nettlesome differences because the currency is specified (or not)
as part of the URL or as a query string.  These differences confound the http client and the formation of request signatures.

We have therefore chosen to deal with this as two separate endpoints.  Please also see account-deposit-history.

*/
func ProbeAccountDepositHistoryByCur(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, currency string) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/deposit/history/"

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

	// 2. This probe has an additional parameter that might be wrong.  Test for these errors if requested.

	// Recall that in this probe the currency parameter is part of the URL. It's not in the queryString.
	var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency. Instead of an error, we would receive status 200 and all the deposits.
		// But that's not an error so don't bother testing it here.

		// 2.2 Request an invalid currency
		invalidCur := "catfood"
		req, err := standardGETReq(credentials, endPoint+invalidCur, "", baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidCur))
		if err != nil {
			fmt.Println("Error with 'currency' param 2.2: ", err)
			return
		}
	}

	// 3. After we've tried all the errors, it's time to build and submit the final correct request.

	// 3.1 Build a request
	req, err := standardGETReq(credentials, endPoint, queryString, baseURL)
	if err != nil {
		fmt.Println("Error building the request 3.1 : ", err)
		return
	}

	// 2.2 We expect a 2xx response
	body, err = TestitAPI2xx(httpClient, req, utils.ExpectedResponseHeaders)
	if err != nil {
		fmt.Println("Error invoking the API 2.2: ", err)
		return
	}
	fmt.Println(body)

	// 2.3 Ensure that the prior response is parsable.
	depositHistories := make([]utils.DepositHistory, 0)
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositHistories)
	if err != nil {
		fmt.Println("Error parsing string into json 2.3: ", err)
		return
	}
}
