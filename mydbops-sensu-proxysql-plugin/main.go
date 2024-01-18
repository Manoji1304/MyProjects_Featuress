/*
Copyright © 2022 Akash N <akash@mydbops.com>

*/
package main

import (
	"fmt"
	"mydbproxysqlchecks/cmd"
	"os"
	"runtime/debug"
)

func recovery() {
	if r := recover(); r != nil {
		fmt.Println("Unknown: ", r)
		debug.PrintStack()
		os.Exit(3)
	}
}

func main() {
	defer recovery()
	cmd.Execute()
}
