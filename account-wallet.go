package main

import (
	"bytes"
	"encoding/json"
	utils "github.com/bostontrader/okcommon"
)

func ProbeAccountWallet(baseURL string, credentialsFile string, makeErrorsCredentials bool) {

	endPoint := "/api/account/v3/wallet"
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
		walletEntries := make([]utils.WalletEntry, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&walletEntries)
		isJSONError(body, err)

		println(string(body))
	}

}
