package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"log"
	"strings"
)

var prefix string = "<SMLK "

//SmlkShell is the shell of SMLKBOT
func SmlkShell(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	if strings.HasPrefix(MsgInfo.Message, prefix) {
		log.Println("Known command: SmlkShell")
		commandstr := strings.Replace(MsgInfo.Message, prefix, "", 1)
		msgArray := strings.Split(commandstr, " ")
		var command = commandMap[msgArray[0]]
		if command != nil {
			go command(MsgInfo, BotConfig, msgArray)
		} else {
			ShellLog(MsgInfo, BotConfig, "notfound")
		}
	}
}
