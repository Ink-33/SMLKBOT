package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/plugins/vtbmusic"
)

//The format for SmlkShell command .
type commandFormat func(*botstruct.MsgInfo, *botstruct.BotConfig, []string)

var commandMap = make(map[string]commandFormat)

func init() {
	commandMap["ping"] = ping
	commandMap["status"] = status
	commandMap["gc"] = gc
	commandMap["reload"] = reload
}
func functionReload() {
	vtbmusic.Load()
}
