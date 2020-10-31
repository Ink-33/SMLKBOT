package music

import (
	"SMLKBOT/data/botstruct"
	"regexp"
	"strings"
)

//VTBMusicClient : VTBMusic点歌Client
type VTBMusicClient struct {
	MusicList []GetVTBMusicListData
}

//AMusicClient : AppleMusic点歌Client
type AMusicClient struct{}

//Client : 点歌Client
type Client interface {
	getMusicDetailandCQCode(int) string
	musicListLen() int
}

//MsgType
type msgType struct {
	//Search Keyword
	content string
	ctype   int8
}

//MsgHandler : The message Handler of music plugin
func MsgHandler(FR *botstruct.FunctionRequest) {
	mt := new(msgType)
	mt.content = ""
	mt.ctype = 0
	if strings.Contains(FR.Message, "CQ:") {
		return
	}
	if strings.HasPrefix(FR.Message, "vtb") {
		vtb(FR, mt)
	}
	if strings.HasPrefix(FR.Message, "苹果点歌") {
		return
	}
	reg := regexp.MustCompile("^[0-9]+$")
	mt.content = strings.Join(reg.FindAllString(FR.Message, 1), "")
	if mt.content != "" {
		if counter != 0 {
			wc := new(waitingChan)
			wc.FunctionRequest = *FR
			wc.isTimeOut = false
			waiting <- wc
		}
	}
}

func vtb(FR *botstruct.FunctionRequest, mt *msgType) {
	if strings.HasPrefix(FR.Message, "vtb点歌") {
		if strings.HasPrefix(FR.Message, "vtb点歌 ") {
			mt.content = strings.Replace(FR.Message, "vtb点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(FR.Message, "vtb点歌", "", 1)
		}
		if mt.content != "" {
			mt.ctype = 1
			VTBMusic(FR, mt)
			return
		}
		mt.ctype = 2
		VTBMusic(FR, mt)
		return
	}

	if strings.HasPrefix(FR.Message, "vtbid点歌") {
		if strings.HasPrefix(FR.Message, "vtbid点歌 ") {
			mt.content = strings.Replace(FR.Message, "vtbid点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(FR.Message, "vtbid点歌", "", 1)
		}
		mt.ctype = 4
		VTBMusic(FR, mt)
		return
	}
	if FR.Message == "vtbhelp" {
		mt.ctype = 5
		mt.content = "Get help"
		VTBMusic(FR, mt)
		return
	}
}
