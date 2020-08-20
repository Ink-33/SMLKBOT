package botstruct

import (
	"github.com/tidwall/gjson"
)

//MsgInfo includes some basic info about a message.
type MsgInfo struct {
	TimeStamp int64
	SenderID  string
	Role      string
	GroupID   string
	Message   string
	MsgType   string
	RobotID   string
	MD5       [16]byte
}

//BotConfig includes OneBot config.
type BotConfig struct {
	MasterID          []gjson.Result
	HTTPAPIAddr       string
	HTTPAPIToken      string
	HTTPAPIPostSecret string
}

/*
Role is the role of message sender.
	RoleLevel => RoleName
		0 => member
		1 => admin
		2 => owner
		3 => master
*/
type Role struct {
	RoleName  string
	RoleLevel int
}
