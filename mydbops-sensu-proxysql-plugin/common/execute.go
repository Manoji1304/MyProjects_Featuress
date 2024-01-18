package common

import (
	"fmt"
	"mydbproxysqlchecks/version"
)

type mydbCheck interface {
	CheckArgs() (int, error)
	ExecuteCheck() (int, error)
}

type Check struct {
	MyCheck mydbCheck
}

func GetCheck(c mydbCheck) *Check {
	return &Check{MyCheck: c}
}

func Alert(status int, message string, mydata interface{}) {
	message = message + " |Version: " + version.Version
	if RootOptions.Json {
		jsonOutput(mydbJsonDesign{ExitStatus: status, Message: message, Output: mydata})
	} else {
		fmt.Println(message)
	}
}
