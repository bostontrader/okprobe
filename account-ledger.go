package main

import (
	"encoding/json"
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"os"
	"strings"
)

func ProbeAccountLedger(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams bool, queryString string) {

	var req *http.Request
	var err error

	// 1. Standard prolog.
	endPoint := "/api/account/v3/ledger"

	// 1.1 Read and parse credentials file
	credentials, err := getCredentials(credentialsFile)
	if err != nil {
		fmt.Println("Error obtaining the credentials 1.1: ", err)
		return
	}

	// 1.2 Obtain an http client
	httpClient := GetHttpClient(baseURL)

	// 1.3 If we want to test header/credentials errors.
	if makeErrorsCredentials {
		err := TestitCredentialsErrors(httpClient, baseURL+endPoint, credentials, utils.ExpectedResponseHeaders)
		if err != nil {
			os.Exit(1)
		}
	}

	// 2. This probe has several additional parameters.  Test for errors if requested.
	// All of the parameters are optional so don't test for their absence.  Only test for bad params.

	// type, before, after, and limit parameters are strings that should parse into integers.  There are subtle differences between the behavior for these things.

	var body string
	var extraExpectedResponseHeaders map[string]string

	if makeErrorsParams {

		var invalidParam string
		var queryString string

		// 2.1 Request an invalid currency
		invalidParam = "catfood"
		queryString = "?currency=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.1: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30031(invalidParam))
		if err != nil {
			fmt.Println("Error with 'currency' param 2.1: ", err)
			return
		}

		// 2.2 Request an invalid type.

		// 2.2.1 If the param cannot parse into an integer, error 500.
		invalidParam = "catfood"
		queryString = "?type=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2.1: ", err)
			return
		}
		body, err = TestitAPI5xx(httpClient, req, utils.ExpectedResponseHeaders, utils.Err500())
		if err != nil {
			fmt.Println("Error with 'type' param 2.2.1: ", err)
			return
		}

		// 2.2.2 If the param can parse into an integer, test that it's one of the chosen ints.
		invalidParam = "666"
		queryString = "?type=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.2.2: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30024("catfood"))
		if err != nil {
			fmt.Println("Error with 'type' param 2.2.2: ", err)
			return
		}

		// 2.3 Request an invalid after.

		// 2.3.1 If the param cannot parse into an integer then error.
		invalidParam = "catfood"
		queryString = "?after=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.3.1: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30025("after parameter format is error"))
		if err != nil {
			fmt.Println("Error with 'after' param 2.3.1: ", err)
			return
		}

		// 2.3.2 If it can parse into any integer then expect success.
		invalidParam = "-1"
		queryString = "?after=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.3.2: ", err)
			return
		}
		extraExpectedResponseHeaders := map[string]string{
			"Strict-Transport-Security": "",
		}
		body, err = TestitAPI2xx(httpClient, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
		if err != nil {
			fmt.Println("Error invoking the API 2.3.2: ", err)
			return
		}

		// 2.3.2.1 Ensure that the prior response is parsable.
		ledgerEntries := make([]utils.LedgerEntry, 0)
		dec := json.NewDecoder(strings.NewReader(body))
		dec.DisallowUnknownFields()
		err = dec.Decode(&ledgerEntries)
		if err != nil {
			fmt.Println("Error parsing string into json 2.3.2.1: ", err)
			return
		}

		// 2.4 Request an invalid before.

		// 2.4.1 If the param cannot parse into an integer then error.
		invalidParam = "catfood"
		queryString = "?before=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.4.1: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30025("before parameter format is error"))
		if err != nil {
			fmt.Println("Error with 'before' param 2.4.1: ", err)
			return
		}

		// 2.4.2 If it can parse into an integer then expect success.
		invalidParam = "-1"
		queryString = "?after=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.5.2: ", err)
			return
		}
		extraExpectedResponseHeaders = map[string]string{
			"Strict-Transport-Security": "",
		}
		body, err = TestitAPI2xx(httpClient, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
		if err != nil {
			fmt.Println("Error invoking the API 2.3.3: ", err)
			return
		}

		// 2.3.3.1 Ensure that the prior response is parsable.
		ledgerEntries = make([]utils.LedgerEntry, 0)
		dec = json.NewDecoder(strings.NewReader(body))
		dec.DisallowUnknownFields()
		err = dec.Decode(&ledgerEntries)
		if err != nil {
			fmt.Println("Error parsing string into json 2.3.3.1: ", err)
			return
		}

		// 2.5 Request an invalid limit.

		// 2.5.1 If the param cannot parse into an integer then error.
		invalidParam = "catfood"
		queryString = "?limit=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.5.1: ", err)
			return
		}
		body, err = TestitAPI4xx(httpClient, req, 400, utils.ExpectedResponseHeaders, utils.Err30025("limit parameter format is error"))
		if err != nil {
			fmt.Println("Error with 'limit' param 2.5.1: ", err)
			return
		}

		// 2.5.2 Limit < 0
		invalidParam = "-1"
		queryString = "?limit=" + invalidParam
		req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
		if err != nil {
			fmt.Println("Error building the request 2.6.2 ", err)
			return
		}
		body, err = TestitAPI5xx(httpClient, req, utils.ExpectedResponseHeaders, utils.Err500())
		if err != nil {
			fmt.Println("Error with 'limit' param 2.5.2: ", err)
			return
		}
	}

	// 3. After we've tried all the errors, it's time to build and submit the final correct request.
	// Since we may optionally feed a query string from the command line, the query string might produce errors.  Deal with it.
	req, err = standardGETReq(credentials, endPoint, queryString, baseURL)
	if err != nil {
		fmt.Println("Error building the request 3: ", err)
		return
	}
	body, err = TestitAPI2xx(httpClient, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	if err != nil {
		fmt.Println("Error invoking the API 3: ", err)
		return
	}
	fmt.Println(body)

	// 3.1 Ensure that the prior response is parsable.
	ledgerEntries := make([]utils.LedgerEntry, 0)
	dec := json.NewDecoder(strings.NewReader(body))
	dec.DisallowUnknownFields()
	err = dec.Decode(&ledgerEntries)
	if err != nil {
		fmt.Println("Error parsing string into json 3.1: ", err)
		return
	}

}