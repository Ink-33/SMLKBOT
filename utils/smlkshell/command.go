package smlkshell

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

// Compile
var date, version, commit = "DevBuild", "DevBuild", "DevBuild"

// IsSCF is the mark to judge whether SMLKBOT is running in SaaS mode.
//	This variable should be set by using -ldflags while building.
var (
	IsSCF  = "no"
	upTime string
)

// RoleHandler : Fetching user's role.
func RoleHandler(fr *botstruct.FunctionRequest) (role *botstruct.Role) {
	role = new(botstruct.Role)
	role.RoleName = "member"
	role.RoleLevel = 0
	for i := range fr.MasterID {
		if fr.SenderID == fr.MasterID[i].String() {
			role.RoleLevel = 3
			role.RoleName = "master"
			return
		}
	}
	if fr.MsgType == "group" {
		role.RoleName = fr.Role
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

// ShellLog : Send execute result to QQ
//	succeed
//	deny
//	disabled
//	nofonud
func ShellLog(fr *botstruct.FunctionRequest, result string) {
	var msgMake string
	switch result {
	case "succeed":
		msgMake = "$SmlkShell> Succeed."
	case "deny":
		msgMake = "$SmlkShell> Permission denied."
	case "disabled":
		msgMake = "$SmlkShell> disabled."
	case "nofonud":
		msgMake = fmt.Sprintf("$SmlkShell> %s: command not found", strings.Replace(fr.Message, prefix, "", 1))
	default:
		msgMake = result
	}
	if msgMake != "" {
		cqfunction.CQSendMsg(fr, msgMake)
	}
}

func ping(fr *botstruct.FunctionRequest, msgArray []string) {
	cost := time.Now().Unix() - fr.TimeStamp
	if len(msgArray) != 1 {
		ShellLog(fr, "nofonud")
		return
	}
	ShellLog(fr, fmt.Sprintf("本次请求耗时:%d秒", cost))
}

func status(fr *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(fr).RoleLevel >= 0 {
		if len(msgArray) != 1 {
			ShellLog(fr, "nofonud")
			return
		}
		msgMake := GetStatus()
		ShellLog(fr, msgMake)
	} else {
		ShellLog(fr, "deny")
	}
}

func gc(fr *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(fr).RoleLevel >= 3 {
		if len(msgArray) != 1 {
			ShellLog(fr, "nofonud")
			return
		}
		runtime.GC()
		ShellLog(fr, "succeed")
	} else {
		ShellLog(fr, "deny")
	}
}

func reload(fr *botstruct.FunctionRequest, msgArray []string) {
	if RoleHandler(fr).RoleLevel == 3 {
		if len(msgArray) != 1 {
			ShellLog(fr, "nofonud")
			return
		}
		cqfunction.ConfigFile = cqfunction.ReadConfig()
		ShellLog(fr, "succeed")
	} else {
		ShellLog(fr, "deny")
	}
}

// GetStatus : Get a string for program status
func GetStatus() string {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	return fmt.Sprintf("$SMLKBOT>\nBuild with: %s\nBuild Arch&OS: %s-%s\nBuild Date: %s\nUptime: %s\nVersion: %s\nCommit: %s\nisSCF: %s\nNumGoroutine: %d\nNumCPU: %d\nNumProcs: %d\nTotalAlloc: %d\nSys(mem): %dBytes\nNumGC: %d\nForceGC: %d\nLsatGC:%s", runtime.Version(), runtime.GOARCH, runtime.GOOS, date, upTime, version, commit, IsSCF, runtime.NumGoroutine(), runtime.NumCPU(), runtime.GOMAXPROCS(0), m.TotalAlloc, m.Sys, m.NumGC, m.NumForcedGC, time.Unix(0, int64(m.LastGC)).Format("2006-01-02 15:04:05"))
}

func init() {
	upTime = time.Now().Format("Mon Jan 2 15:04:05 -0700 MST 2006")
}
