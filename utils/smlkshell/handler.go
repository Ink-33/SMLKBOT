package smlkshell

import (
	"log"
	"strings"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
)

var prefix = "<"

// SmlkShell is the shell of SMLKBOT
func SmlkShell(fr *botstruct.FunctionRequest) {
	if strings.HasPrefix(fr.Message, prefix) {
		msgArray := MsgSplit(fr.Message)
		command := commandMap[msgArray[0]]
		if command != nil {
			log.Println("Known command: SmlkShell")
			go command(fr, msgArray)
		}
	}
}

// MsgSplit : Split massage with " "
func MsgSplit(msg string) (msgArray []string) {
	commandstr := strings.Replace(msg, prefix, "", 1)
	msgArray = strings.Split(commandstr, " ")
	return
}
