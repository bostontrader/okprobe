package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
)

func ProbeAccountDepositAddress(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	// 1. Standard prolog.
	endPoint := "/api/account/v3/deposit/address"

	// 1.1 Read and parse credentials file
	credentials := getCredentials(credentialsFile)

	// 1.2 Obtain an http client
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)

	// 1.3 If we want to test header/credentials errors.
	if makeErrorsCredentials {
		TestitCredentialsErrors(httpClient, url, credentials)
	}

	// 2. This probe has an additional parameter that might be wrong.  Test for these errors if requested.

	// Extraneous parameters appear to be ignored, the full list is returned, and two extra headers appear: Vary and Strict-Transport-Security.  Status = 200.  For example: ?catfood and ?catfood=yum.
	// But this is of minimal importance so don't bother trying to test for this.  We certainly don't care to mimic this behavior in the OKCatbox.
	//var body string
	if makeErrorsParams {
		// 2.1 Don't request any currency
		req, err := standardGETReq(credentials, endPoint, "", baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.1: ", err)
			return
		}
		//_, err = TestitAPI4xxOld(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30023("currency cannot be blank"))
		//if err != nil {
		//fmt.Println("Error with 'currency' param 2.1: ", err)
		//return
		//}
		TestitAPI4xx(httpClient, req, 400, utils.Err30023("currency cannot be blank"))

		// 2.2 Request an invalid currency
		invalidParam := "catfood"
		queryString := "?currency=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2: ", err)
			return
		}
		//_, err = TestitAPI4xxOld(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam))
		//if err != nil {
		//fmt.Println("Error with 'currency' param 2.2: ", err)
		//return
		//}
		TestitAPI4xx(httpClient, req, 400, utils.Err30031(invalidParam))
	}

	// 3. After we've tried all the errors, it's time to build and submit the final correct request.

	// 3.1 Build a request
	req, err := standardGETReq(credentials, endPoint, queryString, baseURL)
	if err != nil {
		fmt.Println("Error building the request 3.1 : ", err)
		return
	}

	// 3.2 We expect a 2xx response
	body1 := TestitAPICore(httpClient, req, 200)

	// 3.3 Ensure that the prior response is parsable.
	depositAddresses := make([]utils.DepositAddress, 0)
	dec := json.NewDecoder(bytes.NewReader(body1))
	dec.DisallowUnknownFields()
	err = dec.Decode(&depositAddresses)
	if err != nil {
		fmt.Println("Error parsing string into json 3.3: ", err)
		return
	}
}
