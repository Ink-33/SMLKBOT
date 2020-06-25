package smlkshell

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
)

//Reload config.
func Reload(MsgInfo *botstruct.MsgInfo, MasterID string) (result bool) {
	result = false
	if MsgInfo.SenderID == MasterID {
		result = true
		return
	}
	return
}

//SmlkShell is the shell of SMLKBOT
func SmlkShell(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {

}

//RoleHandler : Fechting user's role.
func RoleHandler(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) (role string) {
	role = ""
	for _, v := range BotConfig.MasterID {
		if MsgInfo.SenderID == v.String() {
			role = "master"
			return
		}
	}
	if MsgInfo.MsgType == "group" {
		role = MsgInfo.Role
		return
	}
	return
}

//ShellLog : Send execute result to QQ
func ShellLog(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig, result bool) {
	var msgmake string
	if result {
		msgmake = "$SmlkShell> Succeed."
	} else {
		msgmake = "$SmlkShell> Permission denied."
	}
	cqfunction.CQSendMsg(MsgInfo, msgmake, BotConfig)
}
