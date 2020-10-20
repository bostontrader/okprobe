package main

import (
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"strings"
	"time"
)

func ProbeAccountWithdrawal(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType bool, forReal bool, postBody string) {

	endpoint := "/api/account/v3/withdrawal"
	url := baseURL + endpoint
	client := GetHttpClient(baseURL)
	credentials := getCredentialsOld(credentialsFile)

	if makeErrorsCredentials {
		TestitStdPOST(client, url, credentials)
	}

	// Make a call with valid headers but using the wrong credentials.  Wrong credentials will fail first
	// so we don't care about other characteristics of the request.
	if makeErrorsWrongCredentialsType {
		body := "{}"

		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash := timestamp + "POST" + endpoint + body
		encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ := http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		TestitAPI4xx(client, req, 401, utils.Err30012()) // Invalid authority

	}

	if makeErrorsParams {
		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.

		paramTester := ParamTester{
			Client:      client,
			Credentials: credentials,
			Endpoint:    endpoint,
			Url:         url,
		}

		// The API docs say that destination should be an int with a value of 3, 4, or 68 only.
		paramTester.testit(`{}`, utils.Err30025("destination parameter format is error"))
		paramTester.testit(`{"destination":"500"}`, utils.Err30025("fee parameter format is error"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001"}`, utils.Err30023("to_address cannot be blank"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong"}`, utils.Err30023("tradePwd cannot be blank"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "tradePwd":"wrong"}`, utils.Err30023("tradePwd cannot be blank"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong"}`, utils.Err30025("amount parameter format is error"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.1"}`, utils.Err30023("currency cannot be blank"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.1", "currency":"666"}`, utils.Err30031("666"))
		paramTester.testit(`{"destination":"500", "fee":"0.00000001", "to_address":"wrong", "trade_pwd":"wrong", "amount":"0.001", "currency":"XLM"}`, utils.Err34002()) // withdrawal address does not exist
	}

	if forReal {
		// {"result":true,"amount":"0.10000000","from":"6","currency":"BSV","transfer_id":"666666666","to":"1"}
		timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash := timestamp + "POST" + endpoint + postBody
		encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ := http.NewRequest("POST", url, strings.NewReader(postBody))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
		responseBody := TestitAPICore(client, req, 200)

		fmt.Println(string(responseBody))

	}
}
