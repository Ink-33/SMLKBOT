package smlkshell

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

//Compile
var date, version, commit string = "unknown", "unknown", "unknown"

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
	if strings.HasPrefix(MsgInfo.Message, "<SMLK ") {
		if RoleHandler(MsgInfo, BotConfig).RoleLevel >= 1 {
			switch MsgInfo.Message {
			case "<SMLK status":
				m := new(runtime.MemStats)
				runtime.ReadMemStats(m)
				msgMake := fmt.Sprintf("$SMLKBOT>\nBuild with: %s\nBuild Arch&OS: %s-%s\nBuild Date: %s\nVersion: %s\nCommit: %s\nNumGoroutine: %d\nNumCPU: %d\nMemory: %dBytes\nNumGC: %d\nForceGC: %d\nLsatGC:%s", runtime.Version(), runtime.GOARCH, runtime.GOOS, date, version, commit, runtime.NumGoroutine(), runtime.NumCPU(), m.Sys, m.NumForcedGC, m.NumGC, time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"))
				log.Println(msgMake)
				ShellLog(MsgInfo, BotConfig, msgMake)
				break
			case "<SMLK gc":
				if RoleHandler(MsgInfo, BotConfig).RoleLevel == 3 {
					runtime.GC()
					ShellLog(MsgInfo, BotConfig, "succeed")
				} else {
					ShellLog(MsgInfo, BotConfig, "deny")
				}
				break
			default:
				ShellLog(MsgInfo, BotConfig, "notfound")
			}
		} else {
			ShellLog(MsgInfo, BotConfig, "deny")
		}
	}
}

//RoleHandler : Fechting user's role.
func RoleHandler(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) (role *botstruct.Role) {
	role = new(botstruct.Role)
	role.RoleName = "member"
	role.RoleLevel = 0
	for _, v := range BotConfig.MasterID {
		if MsgInfo.SenderID == v.String() {
			role.RoleLevel = 3
			role.RoleName = "master"
			return
		}
	}
	if MsgInfo.MsgType == "group" {
		role.RoleName = MsgInfo.Role
		switch role.RoleName {
		case "owner":
			role.RoleLevel = 2
			return
		case "admin":
			role.RoleLevel = 1
			return
		default:
			return
		}
	}
	return
}

//ShellLog : Send execute result to QQ
func ShellLog(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig, result string) {
	var msgMake string
	switch result {
	case "succeed":
		msgMake = "$SmlkShell> Succeed."
		break
	case "deny":
		msgMake = "$SmlkShell> Permission denied."
		break
	case "notfound":
		msgMake = fmt.Sprintf("$SmlkShell> %s: command not found", strings.Replace(MsgInfo.Message, "<SMLK ", "", 1))
		break
	default:
		msgMake = result
	}
	if msgMake != "" {
		cqfunction.CQSendMsg(MsgInfo, msgMake, BotConfig)
	}
}
