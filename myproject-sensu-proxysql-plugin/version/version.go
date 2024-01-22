package version

import "runtime"

/*
AUTHOR:    AKASH N (Email: akash@mydbops.com)
CONTRIBUTOR:
*/

var (
	//Application Name
	AppName = "mydbproxysqlchecks"

	// Version stores the version of the current build (e.g. 2.0.0)
	Version = "1.0.0"

	// BuildDate stores the timestamp of the build
	// (e.g. 2017-07-31T13:11:15-0700)
	BuildDate string = ""

	// License stores the Copyright Information of the application
	License string = "enterprise"

	// GoVersion stores the version of Go used to build the binary
	GoVersion string = runtime.Version()
)
