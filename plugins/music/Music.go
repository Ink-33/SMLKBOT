package music

import (
	"regexp"
	"strings"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
)

// VTBMusicClient : VTBMusic点歌Client
type VTBMusicClient struct {
	MusicList []GetVTBMusicListData
}

// AMusicClient : AppleMusic点歌Client
type AMusicClient struct{}

// Client : 点歌Client
type Client interface {
	getMusicDetailandCQCode(int) string
	musicListLen() int
}

// MsgType
type msgType struct {
	// Search Keyword
	content string
	ctype   int8
}

// MsgHandler : The message Handler of music plugin
func MsgHandler(fr *botstruct.FunctionRequest) {
	mt := new(msgType)
	mt.content = ""
	mt.ctype = 0
	if strings.Contains(fr.Message, "CQ:") {
		return
	}
	if strings.HasPrefix(fr.Message, "vtb") {
		vtb(fr, mt)
	}
	if strings.HasPrefix(fr.Message, "苹果点歌") {
		return
	}
	reg := regexp.MustCompile("^[0-9]+$")
	mt.content = strings.Join(reg.FindAllString(fr.Message, 1), "")
	if mt.content != "" {
		if counter != 0 {
			wc := new(waitingChan)
			wc.FunctionRequest = *fr
			wc.isTimeOut = false
			waiting <- wc
		}
	}
}

func vtb(fr *botstruct.FunctionRequest, mt *msgType) {
	if strings.HasPrefix(fr.Message, "vtb点歌") {
		if strings.HasPrefix(fr.Message, "vtb点歌 ") {
			mt.content = strings.Replace(fr.Message, "vtb点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(fr.Message, "vtb点歌", "", 1)
		}
		if mt.content != "" {
			mt.ctype = 1
			VTBMusic(fr, mt)
			return
		}
		mt.ctype = 2
		VTBMusic(fr, mt)
		return
	}

	if strings.HasPrefix(fr.Message, "vtbid点歌") {
		if strings.HasPrefix(fr.Message, "vtbid点歌 ") {
			mt.content = strings.Replace(fr.Message, "vtbid点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(fr.Message, "vtbid点歌", "", 1)
		}
		mt.ctype = 4
		VTBMusic(fr, mt)
		return
	}
	if fr.Message == "vtbhelp" {
		mt.ctype = 5
		mt.content = "Get help"
		VTBMusic(fr, mt)
		return
	}
}
