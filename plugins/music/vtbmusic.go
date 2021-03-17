package music

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	help "github.com/Ink-33/SMLKBOT/data/helps"
	txc "github.com/Ink-33/SMLKBOT/plugins/txcloudutils"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

type isHotMusic struct {
	is bool
	/*
		0 -> 未找到歌曲
		1 -> 仅推荐
	*/
	types         int8
	TotalQuantity *int
}

// VTBMusic : The main function of VTBMusic
func VTBMusic(fr *botstruct.FunctionRequest, mt *msgType) {
	client := &VTBMusicClient{}
	ctype := mt.ctype
	var isHot *isHotMusic
	switch ctype {
	case 0:
		break
	case 1:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		go cqfunction.CQSendMsg(fr, "Searching...")
		keywordjson := txc.TenKeywordsExtraction(mt.content, 3)
		keywordStruct := new(txc.KeywordsExtractionRespose)
		err := json.Unmarshal([]byte(keywordjson), keywordStruct)
		if err != nil {
			log.Println(err)
			msgMake := "An unexpected error occurred while fetching data, please check console."
			cqfunction.CQSendMsg(fr, msgMake)
			return
		}
		keywordArray := keywordStruct.Response.Keywords
		if keywordStruct.Response.Error != nil || len(keywordArray) == 0 {
			log.Println("NLP:", keywordStruct.Response.Error, keywordArray)
			list1 := client.GetMusicList(mt.content, "MusicName")
			list2 := client.GetMusicList(mt.content, "VtbName")
			if list1.Total == -1 || list2.Total == -1 {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				cqfunction.CQSendMsg(fr, msgMake)
				return
			}
			ListMsg, ListArray := client.listToMsg(list1, list2)
			if len(ListArray) == 0 {
				list := client.GetHotMusicList()
				ListMsg, ListArray := client.listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				client.sendMsg(fr, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				client.sendMsg(fr, ListMsg, ListArray, mt, isHot)
			}
		} else {
			nlpMsg, nlpArray := client.nlpListToMsg(keywordArray)
			if nlpArray == nil {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				log.Println(msgMake)
				cqfunction.CQSendMsg(fr, msgMake)
				return
			}
			if len(nlpArray) == 0 {
				list := client.GetHotMusicList()
				if list.Total == -1 {
					msgMake := "An unexpected error occurred while fetching data, please check console."
					log.Println(msgMake)
					cqfunction.CQSendMsg(fr, msgMake)
					return
				}
				ListMsg, ListArray := client.listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				client.sendMsg(fr, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				client.sendMsg(fr, nlpMsg, nlpArray, mt, isHot)
			}
		}
	case 2:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command: Get hot music.")
		list := client.GetHotMusicList()
		ListMsg, ListArray := client.listToMsg(list)
		var msgMake string
		if list.Total == -1 {
			msgMake = "An unexpected error occurred while fetching data, please check console."
			log.Println(msgMake)
		} else {
			isHot = &isHotMusic{true, 1, &list.Total}
			client.sendMsg(fr, ListMsg, ListArray, nil, isHot)
		}
	case 3:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		if counter != 0 {
			wc := new(waitingChan)
			wc.FunctionRequest = *fr
			wc.isTimeOut = false
			waiting <- wc
		}
	case 4:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		list := client.GetVTBMusicDetail(mt.content)
		var msgMake string
		switch list.Total {
		case -1:
			msgMake = "An unexpected error occurred while fetching data, please check console."
			log.Println(msgMake)
		case 0:
			msgMake = "[CQ:at,qq=" + fr.SenderID + "]\nid:" + mt.content + "没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
		default:
			client.MusicList = list.Data
			msgMake = client.getMusicDetailandCQCode(1)
		}
		cqfunction.CQSendMsg(fr, msgMake)
	case 5:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		switch fr.MsgType {
		case "private":
			go cqfunction.CQSendPrivateMsg(fr.SenderID, help.VTBMusic, &fr.BotConfig)
		case "group":
			go cqfunction.CQSendGroupMsg(fr.GroupID, help.VTBMusic, &fr.BotConfig)
		}
	}
}

func (e *VTBMusicClient) listToMsg(list ...*VTBMusicList) (listMsg *string, listArray []GetVTBMusicListData) {
	q := make([]string, 0)
	listReturn := make([]GetVTBMusicListData, 0)
	c := 1
	for i := range list {
		for j := range list[i].Data {
			listReturn = append(listReturn, list[i].Data[j])
			t := strconv.Itoa(c) + "," + list[i].Data[j].VocalName + "-" + list[i].Data[j].OriginName
			q = append(q, t)
			c++
		}
	}
	msg := strings.Join(q, "\n")
	return &msg, listReturn
}

func (e *VTBMusicClient) getMusicDetail(index int) (info *VTBMusicInfo) {
	i := new(VTBMusicInfo)
	i.MusicID = e.MusicList[index-1].ID
	i.MusicVocal = strings.ReplaceAll(e.MusicList[index-1].VocalName, ",", "、")
	i.MusicName = e.MusicList[index-1].OriginName
	cdn := e.GetMusicCDN("")
	if strings.Contains(e.MusicList[index-1].CDN, ":") {
		reg := regexp.MustCompile(`(\d+):(\d+):(\d+)`)
		match := reg.FindStringSubmatch(e.MusicList[index-1].CDN)
		i.Cover = cdn.match(match[1]) + e.MusicList[index-1].CoverImg
		i.MusicURL = cdn.match(match[2]) + e.MusicList[index-1].Music
	} else {
		i.MusicCDN = cdn.match(e.MusicList[index-1].CDN)
		i.Cover = i.MusicCDN + e.MusicList[index-1].CoverImg
		i.MusicURL = i.MusicCDN + e.MusicList[index-1].Music
	}
	return i
}

func (e *VTBMusicClient) nlpListToMsg(keywordArray []txc.KeywordsExtractionKeywords) (nlpMsg *string, nlpArray []GetVTBMusicListData) {
	list1 := e.GetMusicList(keywordArray[0].Word, "MusicName")
	list2 := e.GetMusicList(keywordArray[0].Word, "VtbName")

	if list1.Total == -1 || list2.Total == -1 {
		msgMake := "An unexpected error occurred while fetching data, please check console."
		log.Printf("Fail while getting music lists, list1:%d, list2:%d", list1.Total, list2.Total)
		return &msgMake, nil
	}

	nlpArray = make([]GetVTBMusicListData, 0)
	nlpMsgArray := make([]string, 0)

	if list1.Total+list2.Total == 0 {
		msgMake := "nope"
		return &msgMake, nlpArray
	}

	_, ListArray := e.listToMsg(list1, list2)

	for k1 := range keywordArray {
		if len(keywordArray) == 1 {
			nlpArray = ListArray
			break
		}
		switch k1 {
		case 1:
			reg := regexp.MustCompile("(?i)(" + keywordArray[k1].Word + ")")
			for k2 := range ListArray {
				if reg.MatchString(ListArray[k2].VocalName) || reg.MatchString(ListArray[k2].OriginName) {
					nlpArray = append(nlpArray, ListArray[k2])
				}
			}
		default:
			reg := regexp.MustCompile("(?i)(" + keywordArray[k1].Word + ")")
			for k2 := range nlpArray {
				if reg.MatchString(nlpArray[k2].VocalName) || reg.MatchString(nlpArray[k2].OriginName) {
					nlpArray = append(nlpArray, nlpArray[k2])
				}
			}
		}
	}
	for k := range nlpArray {
		tmp := strconv.Itoa(k+1) + "," + nlpArray[k].VocalName + "-" + nlpArray[k].OriginName
		nlpMsgArray = append(nlpMsgArray, tmp)
	}
	msg := strings.Join(nlpMsgArray, "\n")
	return &msg, nlpArray
}

func (e *VTBMusicClient) sendMsg(fr *botstruct.FunctionRequest, listMsg *string, listArray []GetVTBMusicListData, msgType *msgType, isHotMusic *isHotMusic) {
	var msgMake string
	var msgtoGroup string
	lens := len(listArray)
	e.MusicList = listArray
	do := func() {
		counter++
		w := new(waitingChan)
		w.IsNewRequest = true
		w.RequestSenderID = fr.SenderID
		w.FunctionRequest = *fr
		w.isTimeOut = false
		waiting <- w
		var client Client = e
		go waitingFunc(client, fr, &fr.BotConfig)
	}
	if isHotMusic.is {
		if isHotMusic.types != 1 {
			msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》没有在VtbMusic上找到结果。以下是VtbMusic的推荐:\n" + *listMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
			cqfunction.CQSendMsg(fr, msgMake)
			go do()
		} else {
			msgMake = "[CQ:at,qq=" + fr.SenderID + "]\nVTBMusic 当前已收录歌曲 " + strconv.Itoa(*isHotMusic.TotalQuantity) + "首。以下是VtbMusic的推荐:\n" + *listMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
			cqfunction.CQSendMsg(fr, msgMake)
			go do()
		}
	} else {
		switch fr.MsgType {
		case "private":
			switch lens {
			case 1:
				msgMake = e.getMusicDetailandCQCode(1)
			default:
				if lens <= 200 {
					msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *listMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
					go do()
				} else {
					msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试,获取使用帮助请发送vtbhelp"
				}
			}
			go cqfunction.CQSendPrivateMsg(fr.SenderID, msgMake, &fr.BotConfig)
		case "group":
			if lens == 1 {
				msgMake = e.getMusicDetailandCQCode(1)
				go cqfunction.CQSendGroupMsg(fr.GroupID, msgMake, &fr.BotConfig)
				break
			}
			if lens <= 15 {
				msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *listMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				go cqfunction.CQSendGroupMsg(fr.GroupID, msgMake, &fr.BotConfig)
				go do()
				break
			}
			if lens <= 40 {
				msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *listMsg + "\n━━━━━━━━━━━━━━\n请在原群聊发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				msgtoGroup = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到" + strconv.Itoa(lens) + "个结果。为防止打扰到他人，本消息采用私聊推送，请检查私信"
				go cqfunction.CQSendPrivateMsg(fr.SenderID, msgMake, &fr.BotConfig)
				do()
			} else {
				msgMake = "[CQ:at,qq=" + fr.SenderID + "]\n《" + msgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试或私聊BOT获取完整列表,获取使用帮助请发送vtbhelp"
				msgtoGroup = msgMake
			}
			go cqfunction.CQSendGroupMsg(fr.GroupID, msgtoGroup, &fr.BotConfig)
		}
	}
}

func (e *VTBMusicClient) getMusicCode(info *VTBMusicInfo) string {
	return fmt.Sprintf("[CQ:music,type=custom,url=https://vtbmusic.com/song?id=%s,audio=%s,title=%s,content=%s,image=%s]", info.MusicID, info.MusicURL, info.MusicName, info.MusicVocal, info.Cover)
}

func (e *VTBMusicClient) getMusicDetailandCQCode(index int) string {
	return e.getMusicCode(e.getMusicDetail(index))
}

func (e *VTBMusicClient) musicListLen() int {
	return len(e.MusicList)
}
