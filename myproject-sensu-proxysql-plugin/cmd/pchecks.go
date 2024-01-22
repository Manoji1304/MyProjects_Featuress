package cmd

import (
	"fmt"
	"mydbproxysqlchecks/common"
	"mydbproxysqlchecks/internal/constant"
	"mydbproxysqlchecks/internal/plugins/pchecks"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(pchecksCmd)
	pchecksCmd.Flags().StringVarP(&pchecksOptions.CheckType, "type", "t", "status", "ProxySQL check type (one of conns,hg,rules,status,var)")
	pchecksCmd.Flags().StringVarP(&pchecksOptions.VarName, "name", "n", "", "ProxySQL variable name to check")
	pchecksCmd.Flags().StringVarP(&pchecksOptions.IncludeHG, "include-hostgroup", "i", "", "ProxySQL hostgroup(s) to include (only applies to '--type hg' checks, accepts comma-separated list)")
	pchecksCmd.Flags().StringVarP(&pchecksOptions.IgnoreHG, "ignore-hostgroup", "g", "", "ProxySQL hostgroup(s) to ignore (only applies to '--type hg' checks, accepts comma-separated list)")
	pchecksCmd.Flags().StringVar(&pchecksOptions.ToProxy, "to-proxy", "", "Specify proxysql used for proxying (aws,gr,galera,replica), required option if '--type hg'")
	pchecksCmd.Flags().BoolVarP(&pchecksOptions.Lower, "lower", "l", false, "Alert if ProxySQL value are LOWER than defined WARN / CRIT thresholds (only applies to 'var' check type")
	pchecksCmd.Flags().StringVarP(&pchecksOptions.WarnString, "warning", "w", "", "Warning threshold")
	pchecksCmd.Flags().StringVarP(&pchecksOptions.CritString, "critical", "c", "", "Critical threshold")
}

var (
	pchecksOptions pchecks.Variables
	pchecksCmd     = &cobra.Command{
		Use:     constant.PchecksCmdName,
		Example: pchecksExample,
		Args:    cobra.NoArgs,
		Short:   constant.PchecksCmdShortDesc,
		Long:    constant.PchecksCmdLongDesc,
		Run: func(_ *cobra.Command, _ []string) {
			check := common.GetCheck(&pchecksOptions)
			call(check)
		},
	}
)

var pchecksExample = `To Check ProxySQL Uptime
	
	` + fmt.Sprintf("%s %s", rootCmd.Use, constant.PchecksCmdName) + ` --mydbtoken $mydbtoken -u user -p pass -H 127.0.0.1 -P 6032 -t status -w 600 -c 300 --lower

To Check ProxySQL Active_Transactions

	` + fmt.Sprintf("%s %s", rootCmd.Use, constant.PchecksCmdName) + ` --mydbtoken $mydbtoken -u user -p pass -H 127.0.0.1 -P 6032 -t var -n Active_Transactions -w 32 -c 64

To Check ProxySQL Connection Pool Usage

	` + fmt.Sprintf("%s %s", rootCmd.Use, constant.PchecksCmdName) + ` --mydbtoken $mydbtoken -u user -p pass -H 127.0.0.1 -P 6032 -t conns -w 95 -c 98

To Check ProxySQL Hostgroup Availability

	` + fmt.Sprintf("%s %s", rootCmd.Use, constant.PchecksCmdName) + ` --mydbtoken $mydbtoken -u user -p pass -H 127.0.0.1 -P 6032 -t hg -w 1 -c 0 --include-hostgroup 2

To Check ProxySQL Query Rule Configuration

	` + fmt.Sprintf("%s %s", rootCmd.Use, constant.PchecksCmdName) + ` --mydbtoken $mydbtoken -u user -p pass -H 127.0.0.1 -P 6032 -t rules --runtime`
