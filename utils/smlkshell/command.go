package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

//Compile
var date, version, commit string = "DevBuild", "DevBuild", "DevBuild"

//IsSCF is the mark to judge whether SMLKBOT is runing in SaaS mode.
//	This varible should be set by using -ldflags while building.
var IsSCF string = "false"
var upTime string

//SmlkShell is the shell of SMLKBOT
func SmlkShell(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	if strings.HasPrefix(MsgInfo.Message, "<SMLK ") {
		log.Println("Known command: SmlkShell")
		switch MsgInfo.Message {
		case "<SMLK status":
			if RoleHandler(MsgInfo, BotConfig).RoleLevel >= 1 {
				msgMake := GetStatus()
				ShellLog(MsgInfo, BotConfig, msgMake)
			} else {
				ShellLog(MsgInfo, BotConfig, "deny")
			}
			break
		case "<SMLK gc":
			if RoleHandler(MsgInfo, BotConfig).RoleLevel == 3 {
				runtime.GC()
				ShellLog(MsgInfo, BotConfig, "succeed")
			} else {
				ShellLog(MsgInfo, BotConfig, "deny")
			}
			break
		case "<SMLK ping":
			go ping(MsgInfo, BotConfig)
			break
		default:
			ShellLog(MsgInfo, BotConfig, "notfound")
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
	case "disabled":
		msgMake = "$SmlkShell> disabled."
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

func ping(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	cost := time.Now().Unix() - MsgInfo.TimeStamp
	ShellLog(MsgInfo, BotConfig, fmt.Sprintf("本次请求耗时:%d秒", cost))
}

//GetStatus : Get program status
func GetStatus() string {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return fmt.Sprintf("$SMLKBOT>\nBuild with: %s\nBuild Arch&OS: %s-%s\nBuild Date: %s\nUptime: %s\nVersion: %s\nCommit: %s\nisSCF: %s\nNumGoroutine: %d\nNumCPU: %d\nNumProcs: %d\nMemory: %dBytes\nNumGC: %d\nForceGC: %d\nLsatGC:%s", runtime.Version(), runtime.GOARCH, runtime.GOOS, date, upTime, version, commit, IsSCF, runtime.NumGoroutine(), runtime.NumCPU(), runtime.GOMAXPROCS(0), m.Sys, m.NumGC, m.NumForcedGC, time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"))
}

func init() {
	upTime = time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}
