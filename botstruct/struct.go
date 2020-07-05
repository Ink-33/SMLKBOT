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

//Auinfo includes some basic info of a Au number.
type Auinfo struct {
	AuNumber   string
	AuStatus   bool
	AuMsg      string
	AuJumpURL  string
	AuCoverURL string
	AuURL      string
	AuTitle    string
	AuDesp     string
	IsTimeOut  bool
}

//BotConfig includes CQHTTPAPI config.
type BotConfig struct {
	MasterID          []gjson.Result
	HTTPAPIAddr       string
	HTTPAPIToken      string
	HTTPAPIPostSecret string
}

//VTBMusicInfo includes the info of a music.
type VTBMusicInfo struct {
	MusicName  string
	MusicID    string
	MusicVocal string
	Cover      string
	MusicURL   string
	MusicCDN   string
}

//VTBMusicList includes the search result.
type VTBMusicList struct {
	Total int64
	Data  []gjson.Result
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
