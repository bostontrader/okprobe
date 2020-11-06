package main

import (
	"bytes"
	"encoding/json"
	utils "github.com/bostontrader/okcommon"
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
func ProbeAccountDepositHistoryByCur(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool) {

	endPoint := "/api/account/v3/deposit/history/"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	// This probe has an additional parameter that might be wrong.  Test for these errors if requested.
	// Recall that in this probe the currency parameter is part of the URL. It's not in the queryString.
	//var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency. Instead of an error, we would receive status 200 and all the deposits.
		// But that's not an error so don't bother testing it here.

		// 2.2 Request an invalid currency
		invalidCur := "catfood"
		req := buildGETRequest(credentials, endPoint+invalidCur, "", baseURL)

		//_, err = TestitAPI4xxOld(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidCur))
		//if err != nil {
		//fmt.Println("Error with 'currency' param 2.2: ", err)
		//return
		//}
		TestitAPI4xx(httpClient, req, 400, utils.Err30031(invalidCur))
	}

	// This endpoint is a GET and should always work with any credentials.
	if makeErrorsWrongCredentialsType {
		req := buildGETRequest(credentials, endPoint, "", baseURL)
		TestitAPICore(httpClient, req, 200)
	}

	if forReal {
		// Build and execute the request
		req := buildGETRequest(credentials, endPoint, queryString, baseURL)
		body := TestitAPICore(httpClient, req, 200)

		// Ensure that the prior response is parsable.
		depositHistories := make([]utils.DepositHistory, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&depositHistories)
		isJSONError("okprobe:account-ledger.go:ProbeAccountDepositHistoryByCur", body, err)

		println(string(body))
	}

}
