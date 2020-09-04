package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"log"
	"strings"
)

var prefix string = "<SMLK "

//SmlkShell is the shell of SMLKBOT
func SmlkShell(FunctionRequest *botstruct.FunctionRequest) {
	if strings.HasPrefix(FunctionRequest.Message, prefix) {
		log.Println("Known command: SmlkShell")
		commandstr := strings.Replace(FunctionRequest.Message, prefix, "", 1)
		msgArray := strings.Split(commandstr, " ")
		var command = commandMap[msgArray[0]]
		if command != nil {
			go command(FunctionRequest, msgArray)
		} else {
			ShellLog(FunctionRequest, "notfound")
		}
	}
}
