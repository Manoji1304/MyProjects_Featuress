package pchecks

import (
	"database/sql"
	"fmt"
	"mydbproxysqlchecks/internal/constant"
)

func prulesCheck(db *sql.DB) (string, int, error) {
	var (
		srvTable    string = "runtime_mysql_query_rules"
		differCount int
	)

	rulesQuery := fmt.Sprintf("select count(*) from (select * from main.%s where active = 1 EXCEPT select * from disk.mysql_query_rules WHERE active = 1)x", srvTable)
	queryErr := db.QueryRow(rulesQuery).Scan(&differCount)
	if queryErr != nil && queryErr != sql.ErrNoRows {
		return "", constant.CheckStateUnknown, queryErr
	} else if differCount == 0 {
		return "OK: ProxySQL Query Rules: DISK / RUNTIME config matches", constant.CheckStateOK, nil
	} else {
		return "CRITICAL: ProxySQL Query Rules: DISK / RUNTIME config does not match", constant.CheckStateCritical, nil
	}
}
