package biliau2card

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

const biliAuAPIAddr string = "https://www.bilibili.com/audio/music-service-c/web/song/info/h5?sid="
const biliAudioPlayURL string = "https://api.bilibili.com/audio/music-service-c/shareUrl/redirectHttp?songid="
const biliAudioJumpURL string = "https://www.bilibili.com/audio/au"

//GetAuInfo : Get Bilibili Audio info
func GetAuInfo(au string) (Auinfo *botstruct.Auinfo) {
	var ai = new(botstruct.Auinfo)
	reg := regexp.MustCompile("[0-9]+")

	ai.AuNumber = strings.Join(reg.FindAllString(au, 1), "")
	ai.AuURL = biliAudioPlayURL + ai.AuNumber
	ai.AuJumpURL = biliAudioJumpURL + ai.AuNumber

	requestAddr := biliAuAPIAddr + ai.AuNumber
	body := string(cqfunction.GetWebContent(requestAddr)[:])

	ai.AuMsg = gjson.Get(body, "msg").String()
	if ai.AuMsg != "success" {
		ai.AuStatus = false
		return ai
	}
	ai.AuStatus = true
	ai.AuCoverURL = gjson.Get(body, "data.h5Songs.cover_url").String()
	ai.AuTitle = gjson.Get(body, "data.h5Songs.title").String()
	ai.AuDesp = gjson.Get(body, "data.h5Songs.author").String()
	return ai
}
