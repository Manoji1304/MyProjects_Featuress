package cmd

import (
	//c
	// "encoding/json"
	"fmt"
	"mydbproxysqlchecks/common"
	"mydbproxysqlchecks/version"
	"mydbproxysqlchecks/internal/constant"
	"os"
	//c
	// "time"

	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands

var (
	//	rootOptions common.PersistantOptions
	rootCmd = &cobra.Command{
		Use:   constant.RootCmdName,
		Version: fmt.Sprintf("%s, build_date=> %s, licence=> %s, go %s", version.Version, version.BuildDate, version.License, version.GoVersion),
		Args:  cobra.MinimumNArgs(1),
		Short: constant.RootCmdShortDesc,
		Long:  constant.RootCmdLongDesc,
		PersistentPreRun: func(_ *cobra.Command, _ []string) {
			//c
			// status, msg := tokenValidation()
			// if msg != "" {
			// 	common.Alert(status, msg, nil)
			// 	os.Exit(status)
			// }
		},
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(3)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&common.RootOptions.Json, "json", false, "Print output in JSON Format")
	rootCmd.PersistentFlags().BoolVar(&common.RootOptions.JsonPretty, "pretty", false, "Pretty print JSON output (only applicable if --json used)")
	rootCmd.PersistentFlags().StringVarP(&common.RootOptions.DefaultsFile, "defaults-file", "f", "", "ProxySQL defaults file path for login mysql (default: $UserHomeDir/.my.cnf)")
	rootCmd.PersistentFlags().StringVarP(&common.RootOptions.User, "user", "u", "", "ProxySQL username (required if --defaults-file not used)")
	rootCmd.PersistentFlags().StringVarP(&common.RootOptions.Pass, "pass", "p", "", "ProxySQL password (required if --defaults-file not used)")
	rootCmd.PersistentFlags().StringVarP(&common.RootOptions.Host, "host", "H", "127.0.0.1", "ProxySQL host")
	rootCmd.PersistentFlags().Uint64VarP(&common.RootOptions.Port, "port", "P", 6032, "ProxySQL port")
	rootCmd.PersistentFlags().BoolVar(&common.RootOptions.EncryptedCreds, "encrypted-creds", false, "Must use this flag if mydbops encrypted credentials string provided in cli for the checks")
	rootCmd.PersistentFlags().StringVar(&common.RootOptions.MydbToken, "mydbtoken", "", "mydbtoken is a required parameter for checking any resources")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
}
//c
// func tokenValidation() (int, string) {
// 	if common.RootOptions.MydbToken == "" && os.Getenv("MYDB_ALIVE_TOKEN") == "" {
// 		return constant.CheckStateUnknown, "UNKNOWN: --mydbtoken is a required parameter"
// 	} else {
// 		if common.RootOptions.MydbToken == "" && os.Getenv("MYDB_ALIVE_TOKEN") != "" {
// 			common.RootOptions.MydbToken = os.Getenv("MYDB_ALIVE_TOKEN")
// 		}
// 		var ts = make(map[string]int64)
// 		token, err := common.DecryptString(common.RootOptions.MydbToken)
// 		if err != nil {
// 			return constant.CheckStateUnknown, fmt.Sprint(err)
// 		}
// 		err = json.Unmarshal([]byte(token), &ts)
// 		if err != nil {
// 			return constant.CheckStateUnknown, "UNKNOWN: seems wrong mydbtoken, please check"
// 		}
// 		timeNow := time.Now().Unix()
// 		if timeNow > ts["ts"] && timeNow < (ts["ts"]+3600) {
// 			grace := (ts["ts"] + 3600) - timeNow
// 			fmt.Printf("warning => mydbtoken is expired, checks will not work after %d seconds, please renew your token asap | ", grace)
// 		} else if time.Now().Unix() > ts["ts"] {
// 			return constant.CheckStateCritical, "CRITICAL: mydbtoken is expired"
// 		}
// 	}
// 	return 0, ""
// }

func call(c *common.Check) {
	status, err := c.MyCheck.CheckArgs()
	if err != nil {
		common.Alert(status, fmt.Sprint("Unknown: ", err), nil)
		os.Exit(status)
	} else {
		status, err := c.MyCheck.ExecuteCheck()
		if err != nil {
			common.Alert(status, fmt.Sprint("Unknown: ", err), nil)
		}
		os.Exit(status)
	}
}
