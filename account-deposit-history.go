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

The former (this probe) retrieves all deposit history for all currencies.  The later retrieves the same but filters by a specified currency.

These are superficially similar but there are subtle and nettlesome differences because the currency is specified (or not)
as part of the URL or as a query string.  These differences confound the http client and the formation of request signatures.

We have therefore chosen to deal with this as two separate endpoints.  Please also see account-deposit-history-by-cur.

*/
func ProbeAccountDepositHistory(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	endPoint := "/api/account/v3/deposit/history"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	// This endpoint does not take any parameters so there's nothing to do with --makeErrorsParams
	if makeErrorsParams {
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
		isJSONError(body, err)

		println(string(body))
	}

}
