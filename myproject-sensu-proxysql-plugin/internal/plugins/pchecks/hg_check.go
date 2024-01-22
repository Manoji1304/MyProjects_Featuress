package pchecks

import (
	"database/sql"
	"fmt"
	"mydbproxysqlchecks/internal/constant"
	"strings"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.TrimSpace(b) == a {
			return true
		}
	}
	return false
}

func phgCheck(db *sql.DB, options *Variables) (string, int, error) {
	var (
		ignoreHG, includeHG, tableName string
		writerHG                       []string
		readerHG                       []string
	)
	if options.IgnoreHG != "" {
		ignoreHG = fmt.Sprintf("AND sl.hostgroup NOT IN (%s)", options.IgnoreHG)
	} else {
		ignoreHG = ""
	}
	if options.IncludeHG != "" {
		includeHG = fmt.Sprintf("AND sl.hostgroup IN (%s)", options.IncludeHG)
	} else {
		includeHG = ""
	}

	hgQuery := fmt.Sprintf(`SELECT hostgroup_id, hostname, status FROM
	runtime_mysql_servers where status != 'ONLINE' %s %s`, includeHG, ignoreHG)

	queryOut, queryErr := db.Query(hgQuery)

	if queryErr != nil && queryErr != sql.ErrNoRows {
		return "", constant.CheckStateUnknown, queryErr
	}
	type HG struct {
		Hg     string `json:"hg,omitempty"`
		Host   string `json:"host,omitempty"`
		Status string `json:"status,omitempty"`
		Type   string `json:"type,omitempty"`
	}
	var hgs []HG
	for queryOut.Next() {
		var res HG
		err := queryOut.Scan(&res.Hg, &res.Host, &res.Status)
		if err != nil {
			return "", constant.CheckStateUnknown, err
		}
		hgs = append(hgs, res)
	}
	if len(hgs) > 0 {
		switch options.ToProxy {
		case "aws":
			tableName = "mysql_aws_aurora_hostgroups"
		case "gr":
			tableName = "mysql_group_replication_hostgroups"
		case "galera":
			tableName = "mysql_galera_hostgroups"
		case "replica":
			tableName = "mysql_replication_hostgroups"
		}

		if tableName != "" {
			hgInfo, queryErr := db.Query(`SELECT writer_hostgroup, reader_hostgroup FROM ` + tableName)

			if queryErr != nil && queryErr != sql.ErrNoRows {
				return "", constant.CheckStateUnknown, queryErr
			}
			for hgInfo.Next() {
				var rhg, whg string
				err := hgInfo.Scan(&whg, &rhg)
				if err != nil {
					return "", constant.CheckStateUnknown, err
				}
				writerHG = append(writerHG, whg)
				readerHG = append(readerHG, rhg)
			}
		}

		var (
			critMsg string
			warnMsg string
		)
		for _, v := range hgs {
			if stringInSlice(v.Hg, writerHG) {
				v.Type = "writer"
			} else if stringInSlice(v.Hg, readerHG) {
				v.Type = "reader"
			} else {
				v.Type = "unknown"
			}
			vals := fmt.Sprintf("'hg: %s, host: %s, status: %s, type: %s',", v.Hg, v.Host, v.Status, v.Type)
			if v.Status == "SHUNNED" {
				warnMsg += vals
			} else {
				critMsg += vals
			}
		}

		var msgPrefix string = "ProxySQL Hostgroups Status:"
		data = hgs
		if critMsg != "" && warnMsg != "" {
			return fmt.Sprintf("%s: %s CRIT - [%s] WARN - [%s]", constant.CheckMsgCritical, msgPrefix, critMsg, warnMsg), constant.CheckStateCritical, nil
		} else if critMsg != "" {
			return fmt.Sprintf("%s: %s [%s]", constant.CheckMsgCritical, msgPrefix, critMsg), constant.CheckStateCritical, nil
		} else if warnMsg != "" {
			return fmt.Sprintf("%s: %s [%s]", constant.CheckMsgWarning, msgPrefix, warnMsg), constant.CheckStateWarning, nil
		} else {
			return "OK: ProxySQL Hostgroups: ALL OK", constant.CheckStateOK, nil
		}
	} else {
		return "OK: ProxySQL Hostgroups Status: ALL OK", constant.CheckStateOK, nil
	}
}
