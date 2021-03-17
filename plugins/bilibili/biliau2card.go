package bilibili

import (
	"log"
	"regexp"
	"strings"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

// GetAu : Get audio number by regexp.
func GetAu(msg string) (au string) {
	if strings.Contains(msg, "CQ:") {
		return ""
	}

	reg := regexp.MustCompile("(?i)au[0-9]+")

	str := strings.Join(reg.FindAllString(msg, 1), "")
	return str
}

// Au2Card : Handle meassage and send music card.
func Au2Card(fr *botstruct.FunctionRequest) {
	au := GetAu(fr.Message)

	if au != "" {
		log.SetPrefix("BiliAu2Card: ")
		log.Println("Known command:", au)
		AuInfo := GetAuInfo(au)

		if !AuInfo.AuStatus {
			msgMake := "BiliAu2Card: AU" + AuInfo.AuNumber + AuInfo.AuMsg
			switch fr.MsgType {
			case "private":
				go cqfunction.CQSendPrivateMsg(fr.SenderID, msgMake, &fr.BotConfig)
			case "group":
				go cqfunction.CQSendGroupMsg(fr.GroupID, msgMake, &fr.BotConfig)
			}
		} else {
			cqCodeMake := "[CQ:music,type=custom,url=" + AuInfo.AuJumpURL + ",audio=" + AuInfo.AuURL + ",title=" + AuInfo.AuTitle + ",content=" + AuInfo.AuDesp + ",image=" + AuInfo.AuCoverURL + "@180w_180h]"
			switch fr.MsgType {
			case "private":
				go cqfunction.CQSendPrivateMsg(fr.SenderID, cqCodeMake, &fr.BotConfig)
			case "group":
				go cqfunction.CQSendGroupMsg(fr.GroupID, cqCodeMake, &fr.BotConfig)
			}
		}
	}
}
