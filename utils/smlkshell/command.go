package smlkshell

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"fmt"
	"runtime"
	"strings"
	"time"
)

//Compile
var date, version, commit string = "DevBuild", "DevBuild", "DevBuild"

//IsSCF is the mark to judge whether SMLKBOT is runing in SaaS mode.
//	This variable should be set by using -ldflags while building.
var IsSCF string
var upTime string

//RoleHandler : Fetching user's role.
func RoleHandler(FunctionRequest *botstruct.FunctionRequest) (role *botstruct.Role) {
	role = new(botstruct.Role)
	role.RoleName = "member"
	role.RoleLevel = 0
	for _, v := range FunctionRequest.MasterID {
		if FunctionRequest.SenderID == v.String() {
			role.RoleLevel = 3
			role.RoleName = "master"
			return
		}
	}
	if FunctionRequest.MsgType == "group" {
		role.RoleName = FunctionRequest.Role
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
//	succeed
//	deny
//	disabled
//	nofonud
func ShellLog(FunctionRequest *botstruct.FunctionRequest, result string) {
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
	case "nofonud":
		msgMake = fmt.Sprintf("$SmlkShell> %s: command not found", strings.Replace(FunctionRequest.Message, prefix, "", 1))
		break
	default:
		msgMake = result
	}
	if msgMake != "" {
		cqfunction.CQSendMsg(FunctionRequest, msgMake)
	}
}

func ping(FunctionRequest *botstruct.FunctionRequest, msgArray []string) {
	cost := time.Now().Unix() - FunctionRequest.TimeStamp
	if len(msgArray) != 1 {
		ShellLog(FunctionRequest, "nofonud")
		return
	}
	ShellLog(FunctionRequest, fmt.Sprintf("本次请求耗时:%d秒", cost))
}

func status(FunctionRequest *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(FunctionRequest).RoleLevel >= 1 {
		if len(msgArray) != 1 {
			ShellLog(FunctionRequest, "nofonud")
			return
		}
		msgMake := GetStatus()
		ShellLog(FunctionRequest, msgMake)
	} else {
		ShellLog(FunctionRequest, "deny")
	}
}

func gc(FunctionRequest *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(FunctionRequest).RoleLevel >= 3 {
		if len(msgArray) != 1 {
			ShellLog(FunctionRequest, "nofonud")
			return
		}
		runtime.GC()
		ShellLog(FunctionRequest, "succeed")
	} else {
		ShellLog(FunctionRequest, "deny")
	}
}
func reload(FunctionRequest *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(FunctionRequest).RoleLevel == 3 {
		if len(msgArray) != 1 {
			ShellLog(FunctionRequest, "nofonud")
			return
		}
		cqfunction.ConfigFile = cqfunction.ReadConfig()
		functionReload()
		ShellLog(FunctionRequest, "succeed")
	} else {
		ShellLog(FunctionRequest, "deny")
	}
}

//GetStatus : Get a string for program status
func GetStatus() string {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return fmt.Sprintf("$SMLKBOT>\nBuild with: %s\nBuild Arch&OS: %s-%s\nBuild Date: %s\nUptime: %s\nVersion: %s\nCommit: %s\nisSCF: %s\nNumGoroutine: %d\nNumCPU: %d\nNumProcs: %d\nAllocated: %dBytes\nNumGC: %d\nForceGC: %d\nLsatGC:%s", runtime.Version(), runtime.GOARCH, runtime.GOOS, date, upTime, version, commit, IsSCF, runtime.NumGoroutine(), runtime.NumCPU(), runtime.GOMAXPROCS(0), m.Sys, m.NumGC, m.NumForcedGC, time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"))
}

func init() {
	upTime = time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}
