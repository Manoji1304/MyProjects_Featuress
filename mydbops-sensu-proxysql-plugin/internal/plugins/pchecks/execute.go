package pchecks

import (
	"fmt"
	"mydbproxysqlchecks/common"
	"mydbproxysqlchecks/internal/constant"
	"time"

	"github.com/hako/durafmt"
)

var (
	data       interface{}
	pSQLUptime int
)

type varData struct {
	VariableName  string `json:"variable_name,omitempty"`
	VariableValue int    `json:"variable_value"`
}

type optionsInfo struct {
	Warning  string `json:"warning,omitempty"`
	Critical string `json:"critical,omitempty"`
}

func (options *Variables) ExecuteCheck() (int, error) {
	var (
		message  string
		status   int
		checkErr error
		optInfo  *optionsInfo
	)

	db, err := common.ConnectProxySQL("")
	if err != nil {
		return constant.CheckStateUnknown, err
	}

	defer db.Close()

	if options.CheckType == "status" {
		options.VarName = "ProxySQL_Uptime"
		_, status, checkErr = pvarCheck(db, options)
		secsUptime := time.Duration(pSQLUptime) * time.Second
		var prefix string
		if status == constant.CheckStateOK {
			prefix = constant.CheckMsgOK
		} else if status == constant.CheckStateWarning {
			prefix = constant.CheckMsgWarning
		} else if status == constant.CheckStateCritical {
			prefix = constant.CheckMsgCritical
		} else {
			prefix = constant.CheckMsgUnknown
		}
		message = fmt.Sprintf("%s[w: %d,c: %d]: %s = %v", prefix, options.WarnThresh, options.CritThresh, options.VarName, durafmt.Parse(secsUptime))
	} else if options.CheckType == "var" {
		message, status, checkErr = pvarCheck(db, options)
	} else if options.CheckType == "conns" {
		message, status, checkErr = pconnCheck(db, options)
	} else if options.CheckType == "hg" {
		message, status, checkErr = phgCheck(db, options)
	} else if options.CheckType == "rules" {
		message, status, checkErr = prulesCheck(db)
	}

	if checkErr != nil {
		return status, checkErr
	}
	if common.RootOptions.Json {
		optInfo = &optionsInfo{Warning: options.WarnString, Critical: options.CritString}
	}

	common.Alert(status, message, common.MydbJsonData{Args: optInfo, Data: data})
	return status, nil
}
