package pchecks

import (
	"database/sql"
	"fmt"
	"mydbproxysqlchecks/common"
	"mydbproxysqlchecks/internal/constant"
)

func pvarCheck(db *sql.DB, options *Variables) (string, int, error) {
	var (
		pSqlVarValue int
	)
	varQuery := fmt.Sprintf("SELECT variable_value FROM stats.stats_mysql_global where variable_name = '%s';", options.VarName)

	queryErr := db.QueryRow(varQuery).Scan(&pSqlVarValue)
	if queryErr != nil && queryErr != sql.ErrNoRows {
		return "", constant.CheckStateUnknown, queryErr
	}

	if options.CheckType == "status" {
		pSQLUptime = pSqlVarValue
	}

	message := fmt.Sprintf("[w: %d,c: %d]: %s = %d", options.WarnThresh, options.CritThresh, options.VarName, pSqlVarValue)
	if common.RootOptions.Json {
		data = &varData{VariableName: options.VarName, VariableValue: pSqlVarValue}
	}
	if !options.Lower {
		if pSqlVarValue >= options.CritThresh {
			return constant.CheckMsgCritical + message, constant.CheckStateCritical, nil
		} else if pSqlVarValue >= options.WarnThresh {
			return constant.CheckMsgWarning + message, constant.CheckStateWarning, nil
		} else if pSqlVarValue < options.WarnThresh {
			return constant.CheckMsgOK + message, constant.CheckStateOK, nil
		} else {
			return constant.CheckMsgUnknown + message, constant.CheckStateUnknown, nil
		}
	} else {
		if pSqlVarValue <= options.CritThresh {
			return constant.CheckMsgCritical + message, constant.CheckStateCritical, nil
		} else if pSqlVarValue <= options.WarnThresh {
			return constant.CheckMsgWarning + message, constant.CheckStateWarning, nil
		} else if pSqlVarValue > options.CritThresh {
			return constant.CheckMsgOK + message, constant.CheckStateOK, nil
		} else {
			return constant.CheckMsgUnknown + message, constant.CheckStateUnknown, nil
		}
	}

}
