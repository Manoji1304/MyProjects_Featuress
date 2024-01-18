package pchecks

import (
	"fmt"
	"mydbproxysqlchecks/internal/constant"
	"strconv"
)

type Variables struct {
	CheckType  string
	VarName    string
	Lower      bool
	ToProxy    string
	WarnThresh int
	CritThresh int
	WarnString string
	CritString string
	IncludeHG  string
	IgnoreHG   string
}

func (options *Variables) CheckArgs() (int, error) {
	var errors string
	switch options.CheckType {
	case "conns", "hg", "rules", "status", "var":
	default:
		errors = "--type must be one of ('conns','hg','rules','status','var'); "
	}

	if options.CheckType == "hg" && options.ToProxy == "" {
		errors += "--to-proxy is a required parameter for type hg"
	} else if options.CheckType == "hg" && options.ToProxy != "aws" && options.ToProxy != "gr" && options.ToProxy != "galera" && options.ToProxy != "replica" {
		errors += "--to-proxy must be one of ('aws','gr','galera','replica'); "
	}
	if options.CheckType == "status" {
		if options.WarnString == "" {
			options.WarnString = "600"
			options.Lower = true
		}
		if options.CritString == "" {
			options.CritString = "300"
			options.Lower = true
		}
	}
	rsErrorStatement := fmt.Sprintf("You must specify --critical and --warning thresholds for check type %s;", options.CheckType)
	if options.CheckType != "rules" && options.CheckType != "hg" && (options.WarnString == "" || options.CritString == "") {
		return constant.CheckStateUnknown, fmt.Errorf(errors + rsErrorStatement)
	}

	if options.CheckType != "rules" && options.CheckType != "hg" {
		var err error
		options.CritThresh, err = strconv.Atoi(options.CritString)
		if err != nil {
			return constant.CheckStateUnknown, err
		}
		options.WarnThresh, err = strconv.Atoi(options.WarnString)
		if err != nil {
			return constant.CheckStateUnknown, err
		}
	}

	switch {
	case options.WarnThresh > options.CritThresh && !options.Lower:
		errors = errors + fmt.Sprintf("You must specify --critical threshold higher than --warning thresholds for regular %s check;", options.CheckType)
	case options.WarnThresh < options.CritThresh && options.Lower:
		errors = errors + fmt.Sprintf("You must specify --warning threshold higher than --critical thresholds for --lower %s check;", options.CheckType)
	}

	if errors != "" {
		return constant.CheckStateUnknown, fmt.Errorf(errors)
	} else {
		return constant.CheckStateOK, nil
	}
}
