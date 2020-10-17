package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"log"
	"strings"
)

var prefix string = "<"

//SmlkShell is the shell of SMLKBOT
func SmlkShell(FunctionRequest *botstruct.FunctionRequest) {
	if strings.HasPrefix(FunctionRequest.Message, prefix) {
		log.Println("Known command: SmlkShell")
		msgArray := MsgSplit(FunctionRequest.Message)
		var command = commandMap[msgArray[0]]
		if command != nil {
			go command(FunctionRequest, msgArray)
		} else {
			ShellLog(FunctionRequest, "notfound")
		}
	}
}

//MsgSplit : Split massage with " "
func MsgSplit(msg string) (msgArray []string) {
	commandstr := strings.Replace(msg, prefix, "", 1)
	msgArray = strings.Split(commandstr, " ")
	return
}
