package bilibili

import (
	"SMLKBOT/cqfunction"
	"log"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const biliAuAPIAddr string = "https://www.bilibili.com/audio/music-service-c/web/song/info/h5?sid="
const biliAudioPlayURL string = "https://api.bilibili.com/audio/music-service-c/shareUrl/redirectHttp?songid="
const biliAudioJumpURL string = "https://www.bilibili.com/audio/au"

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

//GetAuInfo : Get Bilibili Audio info
func GetAuInfo(au string) (info *Auinfo) {
	var ai = new(Auinfo)
	reg := regexp.MustCompile("[0-9]+")

	ai.AuNumber = strings.Join(reg.FindAllString(au, 1), "")
	ai.AuURL = biliAudioPlayURL + ai.AuNumber
	ai.AuJumpURL = biliAudioJumpURL + ai.AuNumber

	requestAddr := biliAuAPIAddr + ai.AuNumber
	body, err := cqfunction.GetWebContent(requestAddr)
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			ai.IsTimeOut = true
			ai.IsTimeOut = false
			ai.AuMsg = "Time out"
			return
		}
		log.Fatalln(err)
	}
	ai.AuMsg = gjson.GetBytes(body, "msg").String()
	if ai.AuMsg != "success" {
		ai.AuStatus = false
		return ai
	}
	ai.AuStatus = true
	ai.AuCoverURL = gjson.GetBytes(body, "data.h5Songs.cover_url").String()
	ai.AuTitle = gjson.GetBytes(body, "data.h5Songs.title").String()
	ai.AuDesp = gjson.GetBytes(body, "data.h5Songs.author").String()
	return ai
}
