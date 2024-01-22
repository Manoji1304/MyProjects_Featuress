package constant

const (
	CheckStateOK       = 0
	CheckStateWarning  = 1
	CheckStateCritical = 2
	CheckStateUnknown  = 3
	CheckMsgOK         = "OK"
	CheckMsgWarning    = "WARNING"
	CheckMsgCritical   = "CRITICAL"
	CheckMsgUnknown    = "UNKNOWN"
	ShellToUse         = "sh"

	RootCmdName      = "mydbproxysqlchecks"
	RootCmdShortDesc = "Mydbops sensu proxysql checks"
	RootCmdLongDesc  = "Mydbops developed proxysql monitoring plugin"

	PchecksCmdName      = "pchecks"
	PchecksCmdShortDesc = "Mydbops sensu proxysql checks"
	PchecksCmdLongDesc  = `This Sensu plugin alerts based on various check types
('conns','hg','rules','status','var') and print output with status code
between 0 to 3 based on warning and critical threshold`
)
