package vtbmusic

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"SMLKBOT/help"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type msgtype struct {
	content string
	ctype   int
}

var waiting = make(chan *waitingChan, 500)
var counter int = 0

type newRequest struct {
	isNewRequest    bool
	RequestSenderID string
}
type waitingChan struct {
	botstruct.MsgInfo
	isTimeOut bool
	newRequest
}

//VTBMusic : The main function of VTBMusic
func VTBMusic(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	mt := msgHandler(MsgInfo.Message)
	ctype := mt.ctype
	switch ctype {
	case 0:
		break
	case 1:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		list1 := GetVTBMusicList(mt.content)
		list2 := GetVTBVocalList(mt.content)
		listmsg, listarray := listtoMsg(list1, list2)
		sendmsg(MsgInfo, BotConfig, &listmsg, listarray, mt)
		break
	case 2:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command: Get quantity of music.")
		list := GetVTBMusicList(mt.content)
		var msgMake string
		msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\nVTBMusic 当前已收录歌曲 " + strconv.FormatInt(list.Total, 10) + "首。获取使用帮助请发送vtbhelp"
		cqfunction.CQSendMsg(MsgInfo, msgMake, BotConfig)
		break
	case 3:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		if counter != 0 {
			wc := new(waitingChan)
			wc.MsgInfo = *MsgInfo
			wc.isTimeOut = false
			waiting <- wc
		}
		break
	case 4:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		list := GetVTBMusicDetail(mt.content)
		var msgMake string
		if list.Total == 0 {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\nid:" + mt.content + "没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
		} else {
			info := getMusicDetail(list.Data, 1)
			msgMake = "[CQ:music,type=custom,url=https://vtbmusic.com/?song_id=" + info.MusicID + ",audio=" + info.MusicURL + ",title=" + info.MusicName + ",content=" + info.MusicVocal + ",image=" + info.Cover + "]"
		}
		cqfunction.CQSendMsg(MsgInfo, msgMake, BotConfig)
		break
	case 5:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		switch MsgInfo.MsgType {
		case "private":
			go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, help.VTBMusic, BotConfig)
			break
		case "group":
			go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, help.VTBMusic, BotConfig)
			break
		}
		break
	}
}

func msgHandler(msg string) (Msgtype *msgtype) {
	mt := new(msgtype)
	mt.content = ""
	mt.ctype = 0
	if strings.Contains(msg, "CQ:") {
		mt.content = ""
		mt.ctype = 0
		return mt
	}
	if strings.HasPrefix(msg, "vtb点歌") {
		if strings.HasPrefix(msg, "vtb点歌 ") {
			mt.content = strings.Replace(msg, "vtb点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(msg, "vtb点歌", "", 1)
		}
		if mt.content != "" {
			mt.ctype = 1
			return mt
		}
		mt.ctype = 2
		return mt
	}
	reg := regexp.MustCompile("^[0-9]+$")
	mt.content = strings.Join(reg.FindAllString(msg, 1), "")
	if mt.content != "" {
		mt.ctype = 3
		mt.content = msg
		return mt
	}

	if strings.HasPrefix(msg, "vtbid点歌") {
		if strings.HasPrefix(msg, "vtbid点歌 ") {
			mt.content = strings.Replace(msg, "vtbid点歌 ", "", 1)
		} else {
			mt.content = strings.Replace(msg, "vtbid点歌", "", 1)
		}
		mt.ctype = 4
		return mt
	}
	if msg == "vtbhelp" {
		mt.ctype = 5
		mt.content = "Get help"
		return mt
	}
	return mt
}

func listtoMsg(list ...*botstruct.VTBMusicList) (ListMsg string, ListArray []gjson.Result) {
	var q []string
	var listreturn []gjson.Result
	for _, v := range list {
		for i, r := range v.Data {
			listreturn = append(listreturn, r.Get("@this"))
			t := strconv.Itoa(i+1) + "," + r.Get("vocal").String() + "-" + r.Get("name").String()
			q = append(q, t)
			if int64(i+1) == v.Total {
				break
			}
		}
	}

	return strings.Join(q, "\n"), listreturn
}

func waitingFunc(list []gjson.Result, MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	go func(MsgInfo *botstruct.MsgInfo) {
		time.Sleep(60 * time.Second)
		wc := new(waitingChan)
		wc.MsgInfo = *MsgInfo
		wc.isTimeOut = true
		waiting <- wc
	}(MsgInfo)
	for {
		c := <-waiting
		if c.isNewRequest && c.RequestSenderID == MsgInfo.SenderID && c.MD5 != MsgInfo.MD5 {
			counter--
			break
		}
		if c.isTimeOut && c.MD5 == MsgInfo.MD5 {
			counter--
			break
		} else if c.TimeStamp > MsgInfo.TimeStamp && isNumber(c.Message) {
			index, err := strconv.Atoi(c.Message)
			if err != nil {
				log.Fatalln(err)
			}
			if index <= len(list) && index > 0 {
				if c.SenderID == MsgInfo.SenderID && c.MsgType == MsgInfo.MsgType {
					switch c.MsgType {
					case "private":
						info := getMusicDetail(list, index)
						cqCodeMake := "[CQ:music,type=custom,url=https://vtbmusic.com/?song_id=" + info.MusicID + ",audio=" + info.MusicURL + ",title=" + info.MusicName + ",content=" + info.MusicVocal + ",image=" + info.Cover + "]"
						counter--
						go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, BotConfig)
						break
					case "group":
						if c.GroupID == MsgInfo.GroupID {
							info := getMusicDetail(list, index)
							cqCodeMake := "[CQ:music,type=custom,url=https://vtbmusic.com/?song_id=" + info.MusicID + ",audio=" + info.MusicURL + ",title=" + info.MusicName + ",content=" + info.MusicVocal + ",image=" + info.Cover + "]"
							counter--
							go cqfunction.CQSendGroupMsg(c.GroupID, cqCodeMake, BotConfig)
							break
						}
					}
					break
				}
			}
		}
	}
}

func getMusicDetail(list []gjson.Result, index int) (info *botstruct.VTBMusicInfo) {
	i := new(botstruct.VTBMusicInfo)
	i.MusicID = list[index-1].Get("Id").String()
	i.MusicVocal = strings.ReplaceAll(list[index-1].Get("vocal").String(), ",", "、")
	i.MusicName = list[index-1].Get("name").String()
	if strings.Contains(list[index-1].Get("CDN").String(), ":") {
		reg := regexp.MustCompile("(\\d+):(\\d+):(\\d+)")
		match := reg.FindStringSubmatch(list[index-1].Get("CDN").String())
		i.Cover = GetVTBMusicCDN(match[1]) + list[index-1].Get("img").String()
		i.MusicURL = GetVTBMusicCDN(match[2]) + list[index-1].Get("music").String()
	} else {
		i.MusicCDN = GetVTBMusicCDN(list[index-1].Get("CDN").String())
		i.Cover = i.MusicCDN + list[index-1].Get("img").String()
		i.MusicURL = i.MusicCDN + list[index-1].Get("music").String()
	}
	return i
}

func isNumber(str string) bool {
	var result bool = false
	reg := regexp.MustCompile("^[0-9]+$")
	tmp := strings.Join(reg.FindAllString(str, 1), "")
	if tmp != "" {
		result = true
		return result
	}
	return result
}

func sendmsg(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig, listmsg *string, listarray []gjson.Result, Msgtype *msgtype) {
	var msgMake string
	lens := len(listarray)
	if lens == 0 {
		msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + Msgtype.content + "》没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
	} else {
		counter++
		w := new(waitingChan)
		w.isNewRequest = true
		w.RequestSenderID = MsgInfo.SenderID
		w.MsgInfo = *MsgInfo
		w.isTimeOut = false
		waiting <- w
		go waitingFunc(listarray, MsgInfo, BotConfig)
	}
	switch MsgInfo.MsgType {
	case "private":
		if lens != 0 {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + Msgtype.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *listmsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放"
		}
		go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake, BotConfig)
		break
	case "group":
		if lens <= 15 {
			go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgMake, BotConfig)
		} else {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + Msgtype.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *listmsg + "\n━━━━━━━━━━━━━━\n请在原群聊发送歌曲对应序号即可播放"
			msgtoGroup := "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + Msgtype.content + "》共找到" + strconv.Itoa(lens) + "个结果。为防止打扰到他人，本消息采用私聊推送，请检查私信。"
			go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgtoGroup, BotConfig)
			go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake, BotConfig)
		}
		break
	}
}
