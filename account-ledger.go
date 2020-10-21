package main

import (
	"bytes"
	"encoding/json"
	"github.com/bostontrader/okcommon"
)

func ProbeAccountLedger(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	endPoint := "/api/account/v3/ledger"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, baseURL+endPoint, "GET", credentials)
	}

	if makeErrorsParams {

		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.
		paramTester := ParamTester{
			Client:      httpClient,
			Credentials: credentials,
			Endpoint:    endPoint,
			Url:         url,
		}

		// queryString = "" results in 200 and all tx.

		paramTester.GET("?currency=", 401, utils.Err30031(""))
		paramTester.GET("?currency=catfood", 400, utils.Err30031("catfood"))
		// queryString = "?currency=BTC" results in 200 and all tx for BTC.

		// queryString = "?type=" results in 200 and all tx.
		// queryString = "?type=catfood" results in 500.
		paramTester.GET("?type=666", 400, utils.Err30024("Invalid type type"))

		// queryString = "?after=" results in 200 and all tx.
		paramTester.GET("?after=catfood", 400, utils.Err30025("after parameter format is error"))

		// queryString = "?before=" results in 200 and all tx.
		paramTester.GET("?before=catfood", 400, utils.Err30025("before parameter format is error"))

		// queryString = "?limit=" results in 200 and all tx.
		paramTester.GET("?limit=catfood", 400, utils.Err30025("limit parameter format is error"))

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
		ledgerEntries := make([]utils.LedgerEntry, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&ledgerEntries)
		isJSONError(body, err)

		println(string(body))
	}
}
