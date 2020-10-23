package main

import (
	"bytes"
	"encoding/json"
	utils "github.com/bostontrader/okcommon"
)

func ProbeAccountDepositAddress(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	endPoint := "/api/account/v3/deposit/address"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	if makeErrorsParams {

		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.
		paramTester := ParamTester{
			Client:      httpClient,
			Credentials: credentials,
			Endpoint:    endPoint,
			Url:         url,
		}

		paramTester.GET("", 400, utils.Err30023("currency cannot be blank"))
		paramTester.GET("?currency=catfood", 400, utils.Err30031("catfood"))
	}

	// This endpoint is a GET and should always work with any credentials.
	if makeErrorsWrongCredentialsType {
		req := buildGETRequest(credentials, endPoint, "?currency=BTC", baseURL)
		TestitAPICore(httpClient, req, 200)
	}

	if forReal {
		// Build and execute the request
		req := buildGETRequest(credentials, endPoint, queryString, baseURL)
		body := TestitAPICore(httpClient, req, 200)

		// Ensure that the prior response is parsable.
		depositAddresses := make([]utils.DepositAddress, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&depositAddresses)
		isJSONError(body, err)

		println(string(body))
	}

}
