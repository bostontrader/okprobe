package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"os"
	"strings"
)

func ProbeAccountDepositAddress(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/deposit/address"

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

	// Extraneous parameters appear to be ignored, the full list is returned, and two extra headers appear: Vary and Strict-Transport-Security.  Status = 200.  For example: ?catfood and ?catfood=yum.
	// But this is of minimal importance so don't bother trying to test for this.  We certainly don't care to mimic this behavior in the OKCatbox.
	var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency
		req, err := standardGETReq(credentials, endPoint, "", baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.1: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30023("currency cannot be blank"))
		if err != nil {
			fmt.Println("Error with 'currency' param 2.1: ", err)
			return
		}

		// 2.2 Request an invalid currency
		invalidParam := "catfood"
		queryString := "?currency=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam))
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
	depositAddresses := make([]utils.DepositAddress, 0)
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositAddresses)
	if err != nil {
		fmt.Println("Error parsing string into json 3.3: ", err)
		return
	}
}
