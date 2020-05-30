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
)

type msgtype struct {
	content string
	ctype   int
}

var waiting = make(chan *waitingChan, 500)
var counter int = 0

type waitingChan struct {
	botstruct.MsgInfo
	isTimeOut bool
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
		log.Println("Known command:", mt.content, "from:", MsgInfo.SenderID)
		list := GetVTBMusicList(mt.content)
		var msgMake string
		if list.Total == 0 {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
		} else {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》共找到" + strconv.FormatInt(list.Total, 10) + "个结果:\n" + listtoMsg(list) + "\n----------\n发送歌曲对应序号即可播放"
			counter++
			go waitingFunc(list, MsgInfo, BotConfig)
		}
		switch MsgInfo.MsgType {
		case "private":
			go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake, BotConfig)
			break
		case "group":
			if list.Total <= 30 {
				go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgMake, BotConfig)
			} else {
				msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》共找到" + strconv.FormatInt(list.Total, 10) + "个结果:\n" + listtoMsg(list) + "\n----------\n请在原群聊发送歌曲对应序号即可播放"
				msgtoGroup := "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》共找到" + strconv.FormatInt(list.Total, 10) + "个结果。为防止打扰到他人，本消息采用私聊推送，请检查私信。"
				go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgtoGroup, BotConfig)
				go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake, BotConfig)
			}
			break
		}
		break
	case 2:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content, "from:", MsgInfo.SenderID)
		if counter != 0 {
			wc := new(waitingChan)
			wc.MsgInfo = *MsgInfo
			wc.isTimeOut = false
			waiting <- wc
		}
		break
	case 3:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content, "from:", MsgInfo.SenderID)
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
	if strings.Contains(msg, "CQ:rich") {
		mt.content = ""
		mt.ctype = 0
		return mt
	}
	if strings.HasPrefix(msg, "vtb点歌") {
		mt.content = strings.Replace(msg, "vtb点歌", "", 1)
		mt.ctype = 1
		return mt
	}
	reg, err := regexp.Compile("^[0-9]+$")
	if err != nil {
		log.Fatalln(err)
	}
	mt.content = strings.Join(reg.FindAllString(msg, 1), "")
	if mt.content != "" {
		mt.ctype = 2
		mt.content = msg
		return mt
	}
	if msg == "vtbhelp" {
		mt.ctype = 3
		mt.content = "Get help"
		return mt
	}
	return mt
}

func listtoMsg(list *botstruct.VTBMusicList) string {
	var q []string
	for i, r := range list.Data {
		t := strconv.Itoa(i+1) + "," + r.Get("vocal").String() + "-" + r.Get("name").String()
		q = append(q, t)
		if int64(i+1) == list.Total {
			return strings.Join(q, "\n")
		}
	}
	return ""
}

func waitingFunc(list *botstruct.VTBMusicList, MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	go func(MsgInfo *botstruct.MsgInfo) {
		time.Sleep(60 * time.Second)
		wc := new(waitingChan)
		wc.MsgInfo = *MsgInfo
		wc.isTimeOut = true
		waiting <- wc
	}(MsgInfo)
	for {
		c := <-waiting
		if c.isTimeOut && c.MD5 == MsgInfo.MD5 {
			counter--
			break
		} else if c.TimeStamp > MsgInfo.TimeStamp {
			index, err := strconv.Atoi(c.Message)
			if err != nil {
				log.Fatalln(err)
			}
			if int64(index) <= list.Total && int64(index) > 0 {
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

func getMusicDetail(list *botstruct.VTBMusicList, index int) (info *botstruct.VTBMusicInfo) {
	i := new(botstruct.VTBMusicInfo)
	i.MusicID = list.Data[index-1].Get("Id").String()
	i.MusicVocal = list.Data[index-1].Get("vocal").String()
	i.MusicName = list.Data[index-1].Get("name").String()
	if strings.Contains(list.Data[index-1].Get("CDN").String(), ":") {
		reg := regexp.MustCompile("(\\d+):(\\d+):(\\d+)")
		match := reg.FindStringSubmatch(list.Data[index-1].Get("CDN").String())
		i.Cover = GetVTBMusicCDN(match[1]) + list.Data[index-1].Get("img").String()
		i.MusicURL = GetVTBMusicCDN(match[2]) + list.Data[index-1].Get("music").String()
	} else {
		i.MusicCDN = GetVTBMusicCDN(list.Data[index-1].Get("CDN").String())
		i.Cover = i.MusicCDN + list.Data[index-1].Get("img").String()
		i.MusicURL = i.MusicCDN + list.Data[index-1].Get("music").String()
	}
	return i
}
