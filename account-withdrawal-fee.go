package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"os"
	"strings"
)

func ProbeAccountWithdrawalFee(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/withdrawal/fee"

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

	var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency. Instead of an error, we would receive status 200 and all the fees.
		// But that's not an error so don't bother testing it here.

		// 2.2 Request an invalid currency
		invalidParam := "catfood"
		queryString := "?currency=" + invalidParam
		req, err := standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2: ", err)
			return
		}
		body, err = TestitAPI4xxOld(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam))
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

	// 3.2 We expect a 2xx response
	body, err = TestitAPI2xx(httpClient, req, utils.ExpectedResponseHeaders)
	if err != nil {
		fmt.Println("Error invoking the API 3.2: ", err)
		return
	}
	fmt.Println(body)

	// 3.3 Ensure that the prior response is parsable.
	withDrawalFees := make([]utils.WithdrawalFee, 0)
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&withDrawalFees)
	if err != nil {
		fmt.Println("Error parsing string into json 3.3: ", err)
		return
	}

}
