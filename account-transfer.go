package main

import (
	"fmt"
	"github.com/bostontrader/okcommon"
	"net/http"
	"strings"
	"time"
)

func ProbeAccountTransfer(urlBase, keyFile string, makeErrors bool) {

	endpoint := "/api/account/v3/transfer"
	url := urlBase + endpoint
	client := GetClient(urlBase)
	credentials := getCredentials(keyFile)

	if makeErrors {
		TestitStdPOST(client, url, credentials, utils.ExpectedResponseHeaders)

	}

	// Now try a variety of parameter errors.  Be sure to use sufficient credentials. These tests will all require a valid signature.
	body := "{}"

	// The final correct request must have a valid signature, so build one now.
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "POST" + endpoint + body
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders := map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 := Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}

	body = `{"from":""}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}

	body = `{"from":"wrong"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	//{"error_message":"from parameter format is error","code":30025,"error_code":"30025","message":"from parameter format is error"}

	body = `{"from":"666"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	//{"error_message":"to parameter format is error","code":30025,"error_code":"30025","message":"to parameter format is error"}

	body = `{"from":"666", "to":"888"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	//{"error_message":"amount parameter format is error","code":30025,"error_code":"30025","message":"amount parameter format is error"}

	body = `{"from":"666", "to":"888", "amount":"0.1"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"error_message":"currency cannot be blank","code":30023,"error_code":"30023","message":"currency cannot be blank"}	fmt.Println(body1, err)

	body = `{"from":"666", "to":"888", "amount":"0.1", "currency":"xlm"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"error_message":"Invalid type from","code":30024,"error_code":"30024","message":"Invalid type from"}

	body = `{"from":"6", "to":"888", "amount":"0.1", "currency":"xlm"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"error_message":"Invalid type from","code":30024,"error_code":"30024","message":"Invalid type from"}

	body = `{"from":"6", "to":"888", "amount":"0.1", "currency":"xlm"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"error_message":"Invalid type to","code":30024,"error_code":"30024","message":"Invalid type to"}

	body = `{"from":"6", "to":"1", "amount":"0.1", "currency":"xlm"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"error_message":"Insufficient funds","code":34008,"error_code":"34008","message":"Insufficient funds"}

	body = `{"from":"6", "to":"1", "amount":"0.1", "currency":"bsv"}`
	timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash = timestamp + "POST" + endpoint + body
	encoded, _ = utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)

	req, err = http.NewRequest("POST", url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	extraExpectedResponseHeaders = map[string]string{
		"Strict-Transport-Security": "",
	}

	body1 = Testit200(client, req, catMap(utils.ExpectedResponseHeaders, extraExpectedResponseHeaders))
	// {"result":true,"amount":"0.10000000","from":"6","currency":"BSV","transfer_id":"213038068","to":"1"}

	fmt.Println(body1, err)

	//walletEntries := make([]utils.WalletEntry, 0)
	//dec := json.NewDecoder(body)
	//dec.DisallowUnknownFields()
	//err = dec.Decode(&walletEntries)
	//if err != nil {
	//panic(err)
	//}
	//fmt.Println(&walletEntries)
	//fmt.Println(reflect.TypeOf(walletEntries))
}
