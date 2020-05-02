package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

/*
Given a client and a request, that we expect will produce a 4xx error, make the API call and examine the response.  Ensure that
the we get the expected headers, error message, and status code.
*/
func Testit4xx(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
	expectedErrorMessage utils.OKError,
	expectedStatusCode int,
) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != expectedStatusCode {
		fmt.Println("error:\nexpected= ", expectedStatusCode, "\nreceived=", resp.StatusCode)
	}

	if resp.StatusCode == 200 {
	} else {
		var obj utils.OKError
		err = json.Unmarshal(body, &obj)
		if err != nil {
			panic(err)
		}
		if !reflect.DeepEqual(obj, expectedErrorMessage) {
			fmt.Println("error:\nexpected=", expectedErrorMessage, "\nreceived=", obj)
		}
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, expectedResponseHeaders)

	return
}

/*
Given a client, and a request that we expect will produce a 200 response, make the API call and examine the response.  Ensure that
the we get the expected headers and response.
*/
func Testit200(
	client *http.Client,
	req *http.Request,
	expectedResponseHeaders map[string]string,
) io.Reader {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("error:\nexpected= ", 200, "\nreceived=", resp.StatusCode)
	}

	// Look for all of the expected headers in received headers.
	compareHeaders(resp.Header, expectedResponseHeaders)

	// Read the body into a []byte and then create a return a new io.Reader using this []byte.  This enables us to close resp.Body, which we must do, and return an io.Reader which the caller needs in order to Decode JSON.
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(body)
}

// Did we receive all the headers we expected? Did we expect all that we received?
func compareHeaders(respHeaders map[string][]string, expectedResponseHeaders map[string]string) {
	// Look for all of the expected headers in received headers.
	for key, _ := range expectedResponseHeaders {
		_, ok := respHeaders[key]
		if ok {
			// expected and present, cool
		} else {
			fmt.Println("KeyE2:", key, " expected, but not present.")
		}
	}

	// Look for all of the received headers in expected headers.
	for key, _ := range respHeaders {
		_, ok := expectedResponseHeaders[key]
		if ok {
			// received and expected, cool
		} else {
			fmt.Println("KeyR2:", key, " present, but not expected.")
		}
	}
}
