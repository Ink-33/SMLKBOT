package smlkshell

import (
	"SMLKBOT/data/botstruct"
)

//The format for SmlkShell command .
type commandFormat func(*botstruct.FunctionRequest, []string)

var commandMap = make(map[string]commandFormat)

func init() {
	commandMap["ping"] = ping
	commandMap["status"] = status
	commandMap["gc"] = gc
	commandMap["reload"] = reload
}

