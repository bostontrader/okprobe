package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {

	endPtr := flag.String("endpnt", "wallet", "The name of the endpoint to invoke")
	keyFilePtr := flag.String("keyfile", "/path/to/apikeys.json", "The name of a file that contains the API keys")
	urlPtr := flag.String("url", "https://www.okex.com", "The URL of the API")

	if len(os.Args) < 2 {
		flag.Usage()
	}

	flag.Parse()
	fmt.Println("endpnt:", *endPtr)
	fmt.Println("keyfile:", *keyFilePtr)
	fmt.Println("url:", *urlPtr)

	switch *endPtr {
	case "currencies":
		ProbeCurrencies(*urlPtr, *keyFilePtr)
	case "wallet":
		ProbeWallet(*urlPtr, *keyFilePtr)
	case "withdrawal-fee":
		ProbeWithdrawalFee(*urlPtr, *keyFilePtr)

	default:
		fmt.Println("Unknown endpoint ", *endPtr)
	}
}
