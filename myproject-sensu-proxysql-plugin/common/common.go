package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mydbproxysqlchecks/internal/constant"
	"os/exec"
)

func Shell(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(constant.ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	_ = cmd.Process.Release()
	return stdout.String(), stderr.String(), err
}

type mydbJsonDesign struct {
	Type       string      `json:"type"`
	ExitStatus int         `json:"exit_status"`
	Message    string      `json:"message"`
	Output     interface{} `json:"output,omitempty"`
}

type MydbJsonData struct {
	Args interface{} `json:"args,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

func jsonOutput(data mydbJsonDesign) {
	data.Type = "Checks"
	var (
		jsonData []byte
		err      error
	)
	if RootOptions.JsonPretty {
		jsonData, err = json.MarshalIndent(data, "", "\t")
	} else {
		jsonData, err = json.Marshal(data)
	}
	if err != nil {
		fmt.Printf("could not marshal json: %s\n", err)
		return
	}

	fmt.Printf("%s\n", jsonData)
}
