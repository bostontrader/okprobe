package main

import (
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"strings"
	"time"
)

func ProbeAccountTransfer(baseURL string, credentialsFile string, makeErrorsCredentials, makeErrorsParams, makeErrorsWrongCredentialsType bool, forReal bool, postBody string) {

	endpoint := "/api/account/v3/transfer"
	url := baseURL + endpoint
	client := GetHttpClient(baseURL)
	credentials := getCredentials(credentialsFile)

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

		//_, _ = TestitAPI4xxOld(
		//client, req, 401,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30012())
		TestitAPI4xx(client, req, 401, utils.Err30012())

	}

	if makeErrorsParams {
		// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.

		//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}
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

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30025("from parameter format is error"))
		TestitAPI4xx(client, req, 400, utils.Err30025("from parameter format is error"))

		//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}
		body = `{"from":""}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30025("from parameter format is error"))
		TestitAPI4xx(client, req, 400, utils.Err30025("from parameter format is error"))

		//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}
		body = `{"from":"wrong"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30025("from parameter format is error"))
		TestitAPI4xx(client, req, 400, utils.Err30025("from parameter format is error"))

		//{"error_message":"to parameter format is error","code":30025,"error_code":"30025","message":"to parameter format is error"}
		body = `{"from":"666"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30025("to parameter format is error"))
		TestitAPI4xx(client, req, 400, utils.Err30025("to parameter format is error"))

		//{"error_message":"amount parameter format is error","code":30025,"error_code":"30025","message":"amount parameter format is error"}
		body = `{"from":"666", "to":"888"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30025("amount parameter format is error"))
		TestitAPI4xx(client, req, 400, utils.Err30025("amount parameter format is error"))

		// {"error_message":"currency cannot be blank","code":30023,"error_code":"30023","message":"currency cannot be blank"}
		body = `{"from":"666", "to":"888", "amount":"0.1"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30023("currency cannot be blank"))
		TestitAPI4xx(client, req, 400, utils.Err30023("currency cannot be blank"))

		// {"error_message":"Invalid type from","code":30024,"error_code":"30024","message":"Invalid type from"}
		body = `{"from":"666", "to":"888", "amount":"0.1", "currency":"xlm"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30024("Invalid type from"))
		TestitAPI4xx(client, req, 400, utils.Err30024("Invalid type from"))

		// {"error_message":"Invalid type to","code":30024,"error_code":"30024","message":"Invalid type to"}
		body = `{"from":"6", "to":"888", "amount":"0.1", "currency":"xlm"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err30024("Invalid type to"))
		TestitAPI4xx(client, req, 400, utils.Err30024("Invalid type to"))

		// {"error_message":"Insufficient funds","code":34008,"error_code":"34008","message":"Insufficient funds"}
		body = `{"from":"6", "to":"1", "amount":"1000", "currency":"btc"}`
		timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
		prehash = timestamp + "POST" + endpoint + body
		encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

		req, _ = http.NewRequest("POST", url, strings.NewReader(body))
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("OK-ACCESS-KEY", credentials.Key)
		req.Header.Add("OK-ACCESS-SIGN", encoded)
		req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
		req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

		//_, _ = TestitAPI4xxOld(
		//client, req, 400,
		//catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders), utils.Err34008())
		TestitAPI4xx(client, req, 400, utils.Err34008())

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
		body := TestitAPICore(client, req, 200)

		fmt.Println(string(body))

	}
}
