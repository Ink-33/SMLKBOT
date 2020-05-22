package vtbmusic

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type msgtype struct {
	content string
	ctype   int
}

var waiting = make(chan *waitingChan, 2147483647)

type waitingChan struct {
	botstruct.MsgInfo
}

//VTBMusic : The main function of VTBMusic
func VTBMusic(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig) {
	mt := msgHandler(MsgInfo.Message)
	ctype := mt.ctype
	switch ctype {
	case 0:
		log.Println("vtbmusic: Ingore message:", MsgInfo.Message, "from:", MsgInfo.SenderID)
		break
	case 1:
		list := GetVTBMusicList(mt.content)
		var msgMake string
		if list.Total == 0 {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》没有在VtbMusic上找到结果，请更换关键词重试"
		} else {
			msgMake = "[CQ:at,qq=" + MsgInfo.SenderID + "]\n《" + mt.content + "》共找到" + strconv.FormatInt(list.Total, 10) + "个结果:\n" + listtoMsg(list) + "\n----------\n发送歌曲对应序号即可播放"
			go waitingFunc(list, MsgInfo, BotConfig)
		}
		switch MsgInfo.MsgType {
		case "private":
			go cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake, BotConfig)
			break
		case "group":
			go cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgMake, BotConfig)
			break
		}
		break
	case 2:
		wc := new(waitingChan)
		wc.MsgInfo = *MsgInfo
		waiting <- wc
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
	} else {
		reg, err := regexp.Compile("^[\\d]+$")
		if err != nil {
			log.Fatalln(err)
		}
		mt.content = strings.Join(reg.FindAllString(msg, 1), "")
		mt.ctype = 2
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
	for {
		//log.Println("123123123")
		c := <-waiting
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
					//log.Println(cqCodeMake)
					go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, BotConfig)
					break
				case "group":
					if c.GroupID == MsgInfo.GroupID {

						info := getMusicDetail(list, index)
						cqCodeMake := "[CQ:music,type=custom,url=https://vtbmusic.com/?song_id=" + info.MusicID + ",audio=" + info.MusicURL + ",title=" + info.MusicName + ",content=" + info.MusicVocal + ",image=" + info.Cover + "]"
						go cqfunction.CQSendGroupMsg(c.GroupID, cqCodeMake, BotConfig)
						break
					}
				}
				break
			}
		}

	}
}

func getMusicDetail(list *botstruct.VTBMusicList, index int) (info *botstruct.VTBMusicInfo) {
	i := new(botstruct.VTBMusicInfo)
	i.MusicID = list.Data[index-1].Get("Id").String()
	i.MusicVocal = list.Data[index-1].Get("vocal").String()
	i.MusicName = list.Data[index-1].Get("name").String()
	if i.MusicCDN == "11:12:13" {
		i.MusicCDN = "11:12:13"
		i.Cover = "https://santiego.gitee.io/vtb-music-source-img/img/" + list.Data[index-1].Get("img").String()
		i.MusicURL = "https://santiego.gitee.io/vtb-music-source-song/song/" + list.Data[index-1].Get("music").String()
	} else {
		i.MusicCDN = GetVTBMusicCDN(list.Data[index-1].Get("CDN").String())
		i.Cover = i.MusicCDN + list.Data[index-1].Get("img").String()
		i.MusicURL = i.MusicCDN + list.Data[index-1].Get("music").String()
	}

	return i
}
