package main

import (
	"bytes"
	"encoding/json"
	utils "github.com/bostontrader/okcommon"
)

func ProbeAccountCurrencies(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	endPoint := "/api/account/v3/currencies"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	// This endpoint is a GET and takes no queryString.  It should always work with any credentials.
	if makeErrorsWrongCredentialsType {
		req := buildGETRequest(credentials, endPoint, "", baseURL)
		TestitAPICore(httpClient, req, 200)
	}

	// This endpoint does not take any parameters so there's nothing to do with --makeErrorsParams
	if forReal {
		// Build and execute the request
		req := buildGETRequest(credentials, endPoint, "", baseURL)
		body := TestitAPICore(httpClient, req, 200)

		// 2.3 Ensure that the prior response is parsable.
		currencyEntries := make([]utils.CurrenciesEntry, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&currencyEntries)
		isJSONError(body, err)
		println(string(body))
	}
}
