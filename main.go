package main

import (
	"encoding/json"
	"flag"
	"fmt"
	utils "github.com/bostontrader/okcommon"
	"io/ioutil"
	"os"
)

func main() {

	endPtr := flag.String("endpnt", "wallet", "The name of the endpoint to invoke")
	keyFilePtr := flag.String("keyfile", "/path/to/apikeys.json", "The name of a file that contains the API keys")
	makeErrorsPtr := flag.Bool("errors", false, "Invoke the API calls using intentional errors.")
	queryPtr := flag.String("query", "", "The query string to send to the server with GET requests.  For example: '?a=good&b=better'")
	urlPtr := flag.String("url", "https://www.okex.com", "The URL of the API")

	// Args[0] is okprobe
	if len(os.Args) <= 1 {
		flag.Usage()
		return
	}

	flag.Parse()
	fmt.Println("Invoking the Probe with the following configuration:")
	fmt.Println("endpnt:", *endPtr)
	fmt.Println("errors:", *makeErrorsPtr)
	fmt.Println("keyfile:", *keyFilePtr)
	fmt.Println("query:", *queryPtr)
	fmt.Println("url:", *urlPtr)

	switch *endPtr {

	// funding
	case "currencies":
		ProbeCurrencies(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "deposit-address":
		ProbeDepositAddress(*urlPtr, *keyFilePtr, *makeErrorsPtr, *queryPtr)
	case "deposit-history":
		ProbeDepositHistory(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "wallet":
		ProbeWallet(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	case "withdrawal-fee":
		ProbeWithdrawalFee(*urlPtr, *keyFilePtr, *makeErrorsPtr, *queryPtr)

	// spot
	case "accounts":
		ProbeAccounts(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	//case "get-orders":
	//ProbeGetOrders(*urlPtr, *keyFilePtr, *makeErrorsPtr)
	//case "post-orders":
	//ProbePostOrders(*urlPtr, *keyFilePtr, *makeErrorsPtr)

	default:
		fmt.Println("Unknown endpoint ", *endPtr)
	}
}

func getCredentials(keyFile string) utils.Credentials {
	var obj utils.Credentials
	data, err := ioutil.ReadFile(keyFile)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func catMap(a, b map[string]string) map[string]string {
	var n = map[string]string{}
	for k, v := range a {
		n[k] = v
	}
	for k, v := range b {
		n[k] = v
	}

	return n
}
