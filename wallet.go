package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"okprobe/utils"
	"os/user"
	"reflect"
	"time"
)

/*
 signing a message
 using: hmac sha256 + base64
  eg:
    message = Pre_hash function comment
    secretKey = E65791902180E9EF4510DB6A77F6EBAE
  return signed string = TO6uwdqz+31SIPkd4I+9NiZGmVH74dXi+Fd5X0EzzSQ=
*/
func HmacSha256Base64Signer(message string, secretKey string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secretKey))
	_, err := mac.Write([]byte(message))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Each response should have these headers.
var expectedResponseHeaders = map[string]string{
	"X-Ratelimit-Remaining-Second": "",
	"X-Kong-Proxy-Latency":         "",
	"X-Ratelimit-Limit-Second":     "",
	"X-Kong-Upstream-Latency":      "",
	"Set-Cookie":                   "",
	"X-Xss-Protection":             "",
	"X-Frame-Options":              "",
	"Via":                          "",
	"Content-Type":                 "",
	"Pragma":                       "",
	"X-Brokerid":                   "",
	"Content-Length":               "",
	"X-Content-Type-Options":       "",
	"Date":                         "",
	"Cache-Control":                "",
	"Expires":                      "",
}

func testit(
	client *http.Client,
	req *http.Request,
	expectedErrorMessage utils.OKError,
	expectedStatusCode int,
) {

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var obj utils.OKError
	err = json.Unmarshal(body, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}
	if !reflect.DeepEqual(obj, expectedErrorMessage) {
		fmt.Println("error:\nexpected=", expectedErrorMessage, "\nreceived=", obj)
	}
	if resp.StatusCode != expectedStatusCode {
		fmt.Println("error:\nexpected= ", expectedStatusCode, "\nreceived=", resp.StatusCode)
	}

	// look at all expected headers in received headers.
	for key, _ := range expectedResponseHeaders {

		//fmt.Println("KeyE1:", key, "Value:", value)

		_, ok := resp.Header[key]
		if ok {
			// expected and present, cool
		} else {
			fmt.Println("KeyE2:", key, " expected, but not present.")
		}
	}

	// look at all received headers in expected headers.
	for key, _ := range resp.Header {

		//fmt.Println("KeyR1:", key, "Value:", value)

		_, ok := expectedResponseHeaders[key]
		if ok {
			// received and expected, cool
		} else {
			fmt.Println("KeyR2:", key, " present, but not expected.")
		}
	}

	return
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	url := "https://www.okex.com/api/account/v3/wallet"

	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}

	// 1.
	req, _ := http.NewRequest("GET", url, nil)
	testit(client, req, utils.Err30001(), 401) // OK-ACCESS-KEY header is required

	// 2.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	testit(client, req, utils.Err30002(), 400) // OK-ACCESS-SIGN header is required

	// 3.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	testit(client, req, utils.Err30003(), 400) // OK-ACCESS-TIMESTAMP header is required

	// 4.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	testit(client, req, utils.Err30005(), 400) // Invalid OK-ACCESS-TIMESTAMP

	time.Sleep(1 * time.Second) // limit 6/sec

	// 5.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z") // expired
	testit(client, req, utils.Err30008(), 400)                        // Request timestamp expired

	// 6.  Set a good time stamp.  The system time is probably close enough to the server to work.  Maybe try to probe how far off the time can be.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	testit(client, req, utils.Err30006(), 401) // Invalid OK-ACCESS-KEY

	// In order to proceed we need to get real credentials.  Read them from a file.
	type APIKey struct {
		Key        string `json:"api_key"`
		SecretKey  string `json:"api_secret_key"`
		Passphrase string `json:"passphrase"`
	}
	var obj APIKey

	data, err := ioutil.ReadFile(user.HomeDir + "/okex-read.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}

	// 7.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	testit(client, req, utils.Err30004(), 400) // OK-ACCESS-PASSPHRASE header is required

	// 8.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	testit(client, req, utils.Err30015(), 400) // Invalid OK_ACCESS_PASSPHRASE

	// 9.
	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	testit(client, req, utils.Err30013(), 401) // Invalid Sign

	// Now build a signature
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + "/api/account/v3/wallet"
	encoded, _ := HmacSha256Base64Signer(prehash, obj.SecretKey)

	req, _ = http.NewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", obj.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", obj.Passphrase)
	resp, err := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("headers:", resp.Header)
	fmt.Println(string(body))
}
