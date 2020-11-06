package main

import (
	"bytes"
	"encoding/json"
	"github.com/bostontrader/okcommon"
)

func ProbeAccountWithdrawalFee(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	endPoint := "/api/account/v3/withdrawal/fee"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	if makeErrorsParams {
		// 2.1 Don't request any currency. Instead of an error, we would receive status 200 and all the fees.
		// But that's not an error so don't bother testing it here.

		// 2.2 Request an invalid currency
		invalidParam := "catfood"
		queryString := "?currency=" + invalidParam
		req := buildGETRequest(credentials, endPoint, queryString, baseURL)
		TestitAPI4xx(httpClient, req, 400, utils.Err30031(invalidParam))
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
		withDrawalFees := make([]utils.WithdrawalFee, 0)
		dec := json.NewDecoder(bytes.NewReader(body))
		dec.DisallowUnknownFields()
		err := dec.Decode(&withDrawalFees)
		isJSONError("okprobe:account-ledger.go:ProbeAccountWithdrawalFee", body, err)

		println(string(body))
	}

}
