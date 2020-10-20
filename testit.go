package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"
)

type ParamTester struct {
	Client      *http.Client
	Credentials utils.Credentials
	Endpoint    string
	Url         string
}

func (pt *ParamTester) testit(body string, expectedErrorMessage utils.OKError) {

	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "POST" + pt.Endpoint + body
	encoded, _ := utils.HmacSha256Base64Signer(prehash, pt.Credentials.SecretKey)

	req, _ := http.NewRequest("POST", pt.Url, strings.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", pt.Credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", pt.Credentials.Passphrase)

	TestitAPI4xx(pt.Client, req, 400, expectedErrorMessage)

}

// Build an http.Request but terminate the process if we find an error.  Do this to reduce boilerplate error code.
func BuildNewRequest(method string, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Printf("okprobe:testit.go:BuildNewRequest: http.NewRequest error %v\nmethod=%s, url=%s", err, method, url)
		os.Exit(1)
	}
	return req
}

// This is a standard sequence of errors related to credentials to submit to endpoints that are invoked via GET.
func TestitCredentialsErrors(httpClient *http.Client, url string, credentials utils.Credentials) {

	req := BuildNewRequest("GET", url, nil)
	TestitAPI4xx(httpClient, req, 401, utils.Err30001())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	TestitAPI4xx(httpClient, req, 400, utils.Err30002())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	TestitAPI4xx(httpClient, req, 400, utils.Err30003())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	TestitAPI4xx(httpClient, req, 400, utils.Err30005())

	time.Sleep(1 * time.Second) // avoid rate limit

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z") // expired
	TestitAPI4xx(httpClient, req, 400, utils.Err30008())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	TestitAPI4xx(httpClient, req, 401, utils.Err30006())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	TestitAPI4xx(httpClient, req, 400, utils.Err30004())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	TestitAPI4xx(httpClient, req, 400, utils.Err30015())

	req = BuildNewRequest("GET", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)
	TestitAPI4xx(httpClient, req, 401, utils.Err30013())

	return
}

// This is a standard sequence of errors related to credentials to submit to endpoints that are invoked via POST.  There are enough differences between this function and TestitStd to justify the existence of this function.
func TestitStdPOST(client *http.Client, url string, credentials utils.Credentials) {

	req := BuildNewRequest("POST", url, nil)
	TestitAPI4xx(client, req, 401, utils.Err30001()) // OK-ACCESS-KEY header is required

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	TestitAPI4xx(client, req, 400, utils.Err30002()) // OK-ACCESS-SIGN header is required

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	TestitAPI4xx(client, req, 400, utils.Err30003()) // OK-ACCESS-TIMESTAMP header is required

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "invalid")
	TestitAPI4xx(client, req, 400, utils.Err30005()) // Invalid OK-ACCESS-TIMESTAMP

	time.Sleep(1 * time.Second) // limit 6/sec

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", "2020-01-01T01:01:01.000Z") // expired
	TestitAPI4xx(client, req, 400, utils.Err30008())                  // Request timestamp expired

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", "wrong")
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	TestitAPI4xx(client, req, 401, utils.Err30006()) // Invalid OK-ACCESS-KEY

	req = BuildNewRequest("POST", url, nil)
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	TestitAPI4xx(client, req, 400, utils.Err30007()) // Invalid Content_Type, please use the application/json format

	// Send body as nil, []byte(""), []byte(``) gives us a 500 error.

	body := []byte(`{}`)
	req = BuildNewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	TestitAPI4xx(client, req, 400, utils.Err30004()) // OK-ACCESS-PASSPHRASE header is required

	req = BuildNewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", "wrong")
	req.Header.Add("OK-ACCESS-TIMESTAMP", time.Now().UTC().Format("2006-01-02T15:04:05.999Z"))
	req.Header.Add("OK-ACCESS-PASSPHRASE", "wrong")
	TestitAPI4xx(client, req, 400, utils.Err30015()) // Invalid OK_ACCESS_PASSPHRASE

}

/*
Given an http client and a request, we want to make the API call and examine the response.  We want to ensure that
the we get the expected status and error messages.  Read and close the response body and return
said body as a []byte.  In case of error, os.exit(1).

This general process is confounded because status 2xx, 4xx, and 5xx are mostly the same but want to deal with error
messages using different types.  So...

TestitAPICore provides the common functionality...
TestitAPI2xx uses the core and does not expect any error.
TestitAPI4xx uses the core and expects a utils.OKError
TestitAPI5xx uses the core and expects a OK500Error
*/
func TestitAPICore(client *http.Client, req *http.Request, expectedStatusCode int,
) []byte {
	methodName := "okprobe:testit.go:TestitAPICore"
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s: client.Do error %v\n", methodName, err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s: ioutil.ReadAll error %v\n", methodName, err)
		os.Exit(1)
	}

	err = resp.Body.Close()
	if err != nil {
		fmt.Printf("%s: Body.Close() error %v\n", methodName, err)
		os.Exit(1)
	}

	if resp.StatusCode != expectedStatusCode {
		fmt.Printf("%s: Status code error: expected=%d, received=%d\nbody=%s\n", methodName, expectedStatusCode, resp.StatusCode, string(body))
		os.Exit(1)
	}

	return body
}

/* Given an http client and an http request, invoke the request and
verify that the received status code and error message match expectations.  If so return nothing, if not, os.exit(1)
*/
//func TestitAPI200(client *http.Client, req *http.Request) {
//TestitAPICore(client, req, 200)
//}

/* Given an http client, an http request, an expected 4xx status code and error message, invoke the request and
verify that the received status code and error message match expectations.  If so return nothing, if not, os.exit(1)
*/
func TestitAPI4xx(client *http.Client, req *http.Request, expectedStatusCode int, expectedErrorMessage utils.OKError) {
	methodName := "okprobe:testit.go:TestitAPI4xx"
	body := TestitAPICore(client, req, expectedStatusCode)

	var errorMessage utils.OKError
	err := json.Unmarshal(body, &errorMessage)
	if err != nil {
		fmt.Printf("%s: json.Unmarshal error: body=%s\n", methodName, string(body))
		os.Exit(1)
	}

	if !reflect.DeepEqual(errorMessage, expectedErrorMessage) {
		fmt.Printf("%s: Error message compare error: expected=%v, received=%v\n", methodName, expectedErrorMessage, errorMessage)
		os.Exit(1)
	}

	return
}
