package pchecks

import (
	"database/sql"
	"fmt"
	"mydbproxysqlchecks/common"
	"mydbproxysqlchecks/internal/constant"
)

func pconnCheck(db *sql.DB, options *Variables) (string, int, error) {
	var (
		ignoreHG, includeHG string
		srvTable = "runtime_mysql_servers"
	)
	if options.IgnoreHG != "" {
		ignoreHG = fmt.Sprintf("AND s.hostgroup NOT IN (%s)", options.IgnoreHG)
	} else {
		ignoreHG = ""
	}
	if options.IncludeHG != "" {
		includeHG = fmt.Sprintf("AND s.hostgroup IN (%s)", options.IncludeHG)
	} else {
		includeHG = ""
	}

	connQuery := fmt.Sprintf(`SELECT hostgroup_id hg, srv_host srv, port, 
cast((ConnUsed*1.0/max_connections)*100 as int) pct_used FROM main.%s r JOIN 
stats.stats_mysql_connection_pool s ON r.hostgroup_id = s.hostgroup AND srv_host = hostname 
AND srv_port = port %s %s ORDER BY pct_used desc`, srvTable, ignoreHG, includeHG)

	queryOut, queryErr := db.Query(connQuery)
	if queryErr != nil && queryErr != sql.ErrNoRows {
		return "", constant.CheckStateUnknown, queryErr
	}
	type HG struct {
		hg      string
		srv     string
		port    string
		pctUsed int
	}

	var (
		fmtOkOutput   string
		fmtWarnOutput string
		fmtCritOutput string
		fmtUnkOutput  string
	)

	type (
		connTOK   map[string]interface{}
		connTWarn map[string]interface{}
		connTCrit map[string]interface{}
		connTUnk  map[string]interface{}
	)
	var (
		connVOK   []connTOK
		connVWarn []connTWarn
		connVCrit []connTCrit
		connVUnk  []connTUnk
	)

	for queryOut.Next() {
		var res HG
		err := queryOut.Scan(&res.hg, &res.srv, &res.port, &res.pctUsed)
		if err != nil {
			return "", constant.CheckStateUnknown, err
		}
		vals := fmt.Sprintf("'hg: %s, srv: %s, port: %s, pct_used: %d%%',", res.hg, res.srv, res.port, res.pctUsed)
		if res.pctUsed < options.WarnThresh {
			fmtOkOutput = fmtOkOutput + vals
		} else if res.pctUsed >= options.CritThresh {
			fmtCritOutput = fmtCritOutput + vals
		} else if res.pctUsed >= options.WarnThresh {
			fmtWarnOutput = fmtWarnOutput + vals
		} else {
			fmtUnkOutput = fmtUnkOutput + vals
		}
		if common.RootOptions.Json {
			if res.pctUsed < options.WarnThresh {
				d := connTOK{"hg": res.hg, "srv": res.srv, "port": res.port, "pct_used": res.pctUsed}
				connVOK = append(connVOK, d)
			} else if res.pctUsed >= options.CritThresh {
				d := connTCrit{"hg": res.hg, "srv": res.srv, "port": res.port, "pct_used": res.pctUsed}
				connVCrit = append(connVCrit, d)
			} else if res.pctUsed >= options.WarnThresh {
				d := connTWarn{"hg": res.hg, "srv": res.srv, "port": res.port, "pct_used": res.pctUsed}
				connVWarn = append(connVWarn, d)
			} else {
				d := connTUnk{"hg": res.hg, "srv": res.srv, "port": res.port, "pct_used": res.pctUsed}
				connVUnk = append(connVUnk, d)
			}
		}
	}
	var msgPrefix string = "ProxySQL Connections:"
	var msgThres string = fmt.Sprintf("[w: %d,c: %d]", options.WarnThresh, options.CritThresh)
	if fmtUnkOutput != "" {
		data = connVUnk
		return fmt.Sprintf("%s%s: %s [%s]", constant.CheckMsgUnknown, msgThres, msgPrefix, fmtUnkOutput), constant.CheckStateUnknown, nil
	} else if fmtCritOutput != "" {
		data = connVCrit
		return fmt.Sprintf("%s%s: %s [%s]", constant.CheckMsgCritical, msgThres, msgPrefix, fmtCritOutput), constant.CheckStateCritical, nil
	} else if fmtWarnOutput != "" {
		data = connVWarn
		return fmt.Sprintf("%s%s: %s [%s]", constant.CheckMsgWarning, msgThres, msgPrefix, fmtWarnOutput), constant.CheckStateWarning, nil
	} else {
		data = connVOK
		if fmtOkOutput == "" {
			fmtOkOutput = "ALL OK"
		}
		return fmt.Sprintf("%s%s: %s [%s]", constant.CheckMsgOK, msgThres, msgPrefix, fmtOkOutput), constant.CheckStateOK, nil
	}
}
