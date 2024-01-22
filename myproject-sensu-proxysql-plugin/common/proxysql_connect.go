package common

import (
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	//c
	//	"git.heimdall.mydbops.com/mydb/golib/cryptography"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sjmudd/mysql_defaults_file"
)

var RootOptions PersistantOptions

func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && strings.Contains(str, ":")
}

//c
// func DecryptString(encryptedText string) (string, error) {
// 	crypto, err := cryptography.NewCrypto()
// 	if err != nil {
// 		return "", err
// 	}
// 	text, err := crypto.DecryptString(encryptedText)
// 	if err != nil {
// 		return "", err
// 	}
// 	return text, nil
// }

func ConnectProxySQL(dbName string) (*sql.DB, error) {
	var (
		db          *sql.DB
		err         error
		connectType string
	)
	//c
	// if RootOptions.EncryptedCreds && RootOptions.User != "" {
	// 	text, err := DecryptString(RootOptions.User)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to decrypt user string: %s", err)
	// 	}
	// 	RootOptions.User = text
	// }

	// if RootOptions.EncryptedCreds && RootOptions.Pass != "" {
	// 	text, err := DecryptString(RootOptions.Pass)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("unable to decrypt password string: %s", err)
	// 	}
	// 	RootOptions.Pass = text
	// }

	if IsIPv6(strings.TrimSpace(RootOptions.Host)) {
		connectType = fmt.Sprintf("tcp([%s]:%d)/%s", RootOptions.Host, RootOptions.Port, dbName)
	} else {
		connectType = fmt.Sprintf("tcp(%s:%d)/%s", RootOptions.Host, RootOptions.Port, dbName)
	}

	if RootOptions.DefaultsFile != "" {
		db, err = mysql_defaults_file.Open(RootOptions.DefaultsFile, dbName)
	} else if RootOptions.User != "" && RootOptions.Pass != "" {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@%s", RootOptions.User, RootOptions.Pass, connectType))
	} else {
		db, err = mysql_defaults_file.Open("", dbName)
	}
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetConnMaxIdleTime(time.Minute * 1)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	return db, err
}
