package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
)

func ProbeAccountWithdrawalFee(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/withdrawal/fee"

	// 1.1 Read and parse credentials file
	credentials := getCredentials(credentialsFile)

	// 1.2 Obtain an http client
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)

	// 1.3 If we want to test header/credentials errors.
	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "GET", credentials)
	}

	// 2. This probe has an additional parameter that might be wrong.  Test for these errors if requested.

	//var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency. Instead of an error, we would receive status 200 and all the fees.
		// But that's not an error so don't bother testing it here.

		// 2.2 Request an invalid currency
		invalidParam := "catfood"
		queryString := "?currency=" + invalidParam
		req := buildGETRequest(credentials, endPoint, queryString, baseURL)

		//_, err = TestitAPI4xxOld(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam))
		//if err != nil {
		//fmt.Println("Error with 'currency' param 2.2: ", err)
		//return
		//}
		TestitAPI4xx(httpClient, req, 400, utils.Err30031(invalidParam))
	}

	// 3. After we've tried all the errors, it's time to build and submit the final correct request.

	// 3.1 Build a request
	req := buildGETRequest(credentials, endPoint, queryString, baseURL)

	// 3.2 We expect a 2xx response
	body1 := TestitAPICore(httpClient, req, 200)

	// 3.3 Ensure that the prior response is parsable.
	withDrawalFees := make([]utils.WithdrawalFee, 0)
	dec := json.NewDecoder(bytes.NewReader(body1))
	dec.DisallowUnknownFields()
	err := dec.Decode(&withDrawalFees)
	if err != nil {
		fmt.Println("Error parsing string into json 3.3: ", err)
		return
	}

}
