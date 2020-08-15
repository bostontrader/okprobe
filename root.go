package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

var (
	baseURL             string
	okexCredentialsFile string
	tryErrors           bool

	rootCmd = &cobra.Command{
		Use:   "okprobe",
		Short: "Work the OKEx API via the command line.",
		Long:  `Invoke the OKEx API, optionally using a variety of incorrect inputs in order to observe its behavior.`,
	}

	accountCurrenciesCmd = &cobra.Command{
		Use:   "accountCurrencies",
		Short: "Invoke GET /api/account/v3/currencies",
		Long:  `Invoke GET /api/account/v3/currencies`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountWallet ", baseURL)
			ProbeAccountCurrencies(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	accountDepositAddressCmd = &cobra.Command{
		Use:   "accountDepositAddress",
		Short: "Invoke GET /api/account/v3/deposit/address",
		Long:  `Invoke GET /api/account/v3/deposit/address`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountDepositAddress ", baseURL)
			ProbeAccountDepositAddress(baseURL, okexCredentialsFile, tryErrors, "")
		},
	}

	accountDepositHistoryCmd = &cobra.Command{
		Use:   "accountDepositHistory",
		Short: "Invoke GET /api/account/v3/deposit/history",
		Long:  `Invoke GET /api/account/v3/deposit/history`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountDepositHistory ", baseURL)
			ProbeAccountDepositHistory(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	accountLedgerCmd = &cobra.Command{
		Use:   "accountLedger",
		Short: "Invoke GET /api/account/v3/ledger",
		Long:  `Invoke GET /api/account/v3/ledger`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountLedger ", baseURL)
			ProbeAccountLedger(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	accountWalletCmd = &cobra.Command{
		Use:   "accountWallet",
		Short: "Invoke GET /api/account/v3/wallet",
		Long:  `Invoke GET /api/account/v3/wallet`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountWallet ", baseURL)
			ProbeAccountWallet(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	accountWithdrawalFeeCmd = &cobra.Command{
		Use:   "accountWithdrawalFee",
		Short: "Invoke GET /api/account/v3/withdrawal/fee",
		Long:  `Invoke GET /api/account/v3/withdrawal/fee`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET accountWithdrawalFee ", baseURL)
			ProbeAccountWithdrawalFee(baseURL, okexCredentialsFile, tryErrors, "")
		},
	}

	spotAccountsCmd = &cobra.Command{
		Use:   "spotAccounts",
		Short: "Invoke GET /api/spot/v3/accounts",
		Long:  `Invoke GET /api/spot/v3/accounts`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking GET spotAccounts ", baseURL)
			ProbeSpotAccounts(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	accountTransferCmd = &cobra.Command{
		Use:   "accountTransfer",
		Short: "Invoke POST /api/account/v3/transfer",
		Long:  `Invoke POST /api/account/v3/transfer`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Invoking POST accountTransfer ", baseURL)
			ProbeAccountWallet(baseURL, okexCredentialsFile, tryErrors)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of OKProbe",
		Long:  `All software has versions. This is OKProbe's`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("OKProbe v0.1 -- HEAD")
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	//cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&baseURL, "url", "https://www.okex.com", "The base url of the OKEx server to use.")
	rootCmd.PersistentFlags().StringVar(&okexCredentialsFile, "credentials-file", "/path/to/mycredentials.json", "A file name for a file that contains the OKEx credentials to use.")
	rootCmd.PersistentFlags().BoolVar(&tryErrors, "errors", false, "When invoking this API first send a variety of known errors and test that we receive the expected responses. (default false)")

	rootCmd.AddCommand(accountCurrenciesCmd)
	rootCmd.AddCommand(accountDepositAddressCmd)
	rootCmd.AddCommand(accountDepositHistoryCmd)
	rootCmd.AddCommand(accountLedgerCmd)
	rootCmd.AddCommand(accountTransferCmd)
	rootCmd.AddCommand(accountWalletCmd)
	rootCmd.AddCommand(accountWithdrawalFeeCmd)
	rootCmd.AddCommand(spotAccountsCmd)

	rootCmd.AddCommand(versionCmd)

}
