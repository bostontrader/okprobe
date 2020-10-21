package main

import (
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"strings"
	"time"
)

func ProbeAccountWithdrawal(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType bool, forReal bool, postBody string) {

	endPoint := "/api/account/v3/withdrawal"
	url := baseURL + endPoint
	httpClient := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

	if makeErrorsCredentials {
		TestitCredentialsHeadersErrors(httpClient, url, "POST", credentials)
	}

	// Make a call with valid headers but using the wrong credentials.  Wrong credentials will fail first
	// so we don't care about other characteristics of the request.
	if makeErrorsWrongCredentialsType {
		body := "{}"

		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash := timestamp + "POST" + endPoint + body
		encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ := http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		TestitAPI4xx(httpClient, req, 401, utils.Err30012()) // Invalid authority
	}

	if makeErrorsParams {
		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.

		paramTester := ParamTester{
			Client:      httpClient,
			Credentials: credentials,
			Endpoint:    endPoint,
			Url:         url,
		}

		// The API docs say that destination should be an int with a value of 3, 4, or 68 only.
		paramTester.POST(`{}`, utils.Err30025("destination parameter format is error"))
		paramTester.POST(`{"destination":"500"}`, utils.Err30025("fee parameter format is error"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001"}`, utils.Err30023("to_address cannot be blank"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong"}`, utils.Err30023("tradePwd cannot be blank"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "tradePwd":"wrong"}`, utils.Err30023("tradePwd cannot be blank"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong"}`, utils.Err30025("amount parameter format is error"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.1"}`, utils.Err30023("currency cannot be blank"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.1", "currency":"666"}`, utils.Err30031("666"))
		paramTester.POST(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.001", "currency":"XLM"}`, utils.Err34002()) // withdrawal address does not exist
	}

	if forReal {
		// Build and execute the request
		req := buildPOSTRequest(credentials, endPoint, postBody, baseURL)
		body := TestitAPICore(httpClient, req, 200)

		// Ensure that the prior response is parsable.
		fmt.Println(string(body))
	}
}
