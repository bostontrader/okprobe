package main

import (
	"encoding/json"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// queryString should have the initial ? mark, if present
func standardGETReq(credentials utils.Credentials, endPoint string, queryString, urlBase string) (*http.Request, error) {
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.999Z")
	prehash := timestamp + "GET" + endPoint + queryString
	encoded, _ := utils.HmacSha256Base64Signer(prehash, credentials.SecretKey)
	req, err := http.NewRequest("GET", urlBase+endPoint+queryString, nil)
	if err != nil {
		fmt.Println("Error building NewRequest ", err)
		return nil, err
	}
	req.Header.Add("OK-ACCESS-KEY", credentials.Key)
	req.Header.Add("OK-ACCESS-SIGN", encoded)
	req.Header.Add("OK-ACCESS-TIMESTAMP", timestamp)
	req.Header.Add("OK-ACCESS-PASSPHRASE", credentials.Passphrase)

	return req, nil
}

/* Given a file name, attempt to read and parse a credentials file.  os.exit(1) on any error. */
func getCredentials(keyFile string) utils.Credentials {
	methodName := "okprobe:main.go:getCredentials"
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		fmt.Printf("%s: Error reading file %s, error=%v\n", methodName, keyFile, err)
		os.Exit(1)
	}

	var obj utils.Credentials
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Printf("%s: Error parsing credentials file, error=%v\ndata=%s", keyFile, err, string(data))
		os.Exit(1)
	}

	return obj
}

func main() {
	err := Execute()
	if err != nil {
		fmt.Printf("okprobe:main.go:main: Execute error %v\n", err)
		os.Exit(1)
	}
}
