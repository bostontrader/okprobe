package main

import (
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"strings"
	"time"
)

func ProbeAccountTransfer(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType bool, forReal bool, postBody string) {

	endPoint := "/api/account/v3/transfer"
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
		TestitAPI4xx(httpClient, req, 401, utils.Err30012())
	}

	if makeErrorsParams {
		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.

		paramTester := ParamTester{
			Client:      httpClient,
			Credentials: credentials,
			Endpoint:    endPoint,
			Url:         url,
		}
		paramTester.POST(`{}`, utils.Err30025("from parameter format is error"))
		paramTester.POST(`{"from":""}`, utils.Err30025("from parameter format is error"))
		paramTester.POST(`{"from":"wrong"}`, utils.Err30025("from parameter format is error"))
		paramTester.POST(`{"from":"666"}`, utils.Err30025("to parameter format is error"))
		paramTester.POST(`{"from":"666"}`, utils.Err30025("to parameter format is error"))

		time.Sleep(1 * time.Second) // rate limit throttle

		paramTester.POST(`{"from":"666", "to":"888"}`, utils.Err30025("amount parameter format is error"))
		paramTester.POST(`{"from":"666", "to":"888", "amount":"0.1"}`, utils.Err30023("currency cannot be blank"))
		paramTester.POST(`{"from":"666", "to":"888", "amount":"0.1", "currency":"xlm"}`, utils.Err30024("Invalid type from"))
		time.Sleep(1 * time.Second) // rate limit throttle

		paramTester.POST(`{"from":"6", "to":"888", "amount":"0.1", "currency":"xlm"}`, utils.Err30024("Invalid type to"))
		paramTester.POST(`{"from":"6", "to":"1", "amount":"10000", "currency":"btc"}`, utils.Err34008())

		// {"error_message":"Insufficient funds","code":34008,"error_code":"34008","message":"Insufficient funds"}
		//body = `{"from":"6", "to":"1", "amount":"1000", "currency":"btc"}`

	}

	if forReal {
		// Build and execute the request
		req := buildPOSTRequest(credentials, endPoint, postBody, baseURL)
		body := TestitAPICore(httpClient, req, 200)

		// Ensure that the prior response is parsable.
		fmt.Println(string(body))
	}
}
