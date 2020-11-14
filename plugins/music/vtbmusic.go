package music

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/data/helps"
	txc "SMLKBOT/plugins/txcloudutils"
	"SMLKBOT/utils/cqfunction"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
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

//VTBMusic : The main function of VTBMusic
func VTBMusic(FunctionRequest *botstruct.FunctionRequest, mt *msgType) {
	client := &VTBMusicClient{}
	ctype := mt.ctype
	var isHot *isHotMusic
	switch ctype {
	case 0:
		break
	case 1:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		go cqfunction.CQSendMsg(FunctionRequest, "Searching...")
		keywordjson := txc.TenKeywordsExtraction(mt.content, 3)
		keywordStruct := new(txc.KeywordsExtractionRespose)
		err := json.Unmarshal([]byte(keywordjson), keywordStruct)
		if err != nil {
			log.Println(err)
			msgMake := "An unexpected error occurred while fetching data, please check console."
			cqfunction.CQSendMsg(FunctionRequest, msgMake)
			return
		}
		keywordArray := keywordStruct.Response.Keywords
		if keywordStruct.Response.Error != nil || len(keywordArray) == 0 {
			log.Println("NLP:", keywordStruct.Response.Error, keywordArray)
			list1 := client.GetMusicList(mt.content, "MusicName")
			list2 := client.GetMusicList(mt.content, "VtbName")
			if list1.Total == -1 || list2.Total == -1 {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				cqfunction.CQSendMsg(FunctionRequest, msgMake)
				return
			}
			ListMsg, ListArray := client.listToMsg(list1, list2)
			if len(ListArray) == 0 {
				list := client.GetHotMusicList()
				ListMsg, ListArray := client.listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				client.sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				client.sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			}
		} else {
			nlpMsg, nlpArray := client.nlpListToMsg(keywordArray)
			if nlpArray == nil {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				log.Println(msgMake)
				cqfunction.CQSendMsg(FunctionRequest, msgMake)
				return
			}
			if len(nlpArray) == 0 {
				list := client.GetHotMusicList()
				if list.Total == -1 {
					msgMake := "An unexpected error occurred while fetching data, please check console."
					log.Println(msgMake)
					cqfunction.CQSendMsg(FunctionRequest, msgMake)
					return
				}
				ListMsg, ListArray := client.listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				client.sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				client.sendMsg(FunctionRequest, nlpMsg, nlpArray, mt, isHot)
			}
		}
		break
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
			client.sendMsg(FunctionRequest, ListMsg, ListArray, nil, isHot)
		}
		break
	case 3:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		if counter != 0 {
			wc := new(waitingChan)
			wc.FunctionRequest = *FunctionRequest
			wc.isTimeOut = false
			waiting <- wc
		}
		break
	case 4:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		list := client.GetVTBMusicDetail(mt.content)
		var msgMake string
		if list.Total == -1 {
			msgMake = "An unexpected error occurred while fetching data, please check console."
			log.Println(msgMake)
		} else if list.Total == 0 {
			msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\nid:" + mt.content + "没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
		} else {
			client.MusicList = list.Data
			msgMake = client.getMusicDetailandCQCode(1)
		}
		cqfunction.CQSendMsg(FunctionRequest, msgMake)
		break
	case 5:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		switch FunctionRequest.MsgType {
		case "private":
			go cqfunction.CQSendPrivateMsg(FunctionRequest.SenderID, help.VTBMusic, &FunctionRequest.BotConfig)
			break
		case "group":
			go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, help.VTBMusic, &FunctionRequest.BotConfig)
			break
		}
		break
	}
}

func (e *VTBMusicClient) listToMsg(list ...*VTBMusicList) (ListMsg *string, ListArray []GetVTBMusicListData) {
	var q = make([]string, 0)
	var listReturn = make([]GetVTBMusicListData, 0)
	var c int = 1
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
		reg := regexp.MustCompile("(\\d+):(\\d+):(\\d+)")
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

func (e *VTBMusicClient) nlpListToMsg(keywordArray []txc.KeywordsExtractionKeywords) (NLPMsg *string, NLPArray []GetVTBMusicListData) {
	list1 := e.GetMusicList(keywordArray[0].Word, "MusicName")
	list2 := e.GetMusicList(keywordArray[0].Word, "VtbName")

	if list1.Total == -1 || list2.Total == -1 {
		msgMake := "An unexpected error occurred while fetching data, please check console."
		log.Printf("Fail while getting music lists, list1:%d, list2:%d", list1.Total, list2.Total)
		return &msgMake, nil
	}

	var nlpArray = make([]GetVTBMusicListData, 0)
	var nlpMsgArray = make([]string, 0)

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

func (e *VTBMusicClient) sendMsg(FunctionRequest *botstruct.FunctionRequest, ListMsg *string, ListArray []GetVTBMusicListData, MsgType *msgType, isHotMusic *isHotMusic) {
	var msgMake string
	var msgtoGroup string
	lens := len(ListArray)
	e.MusicList = ListArray
	do := func() {
		counter++
		w := new(waitingChan)
		w.isNewRequest = true
		w.RequestSenderID = FunctionRequest.SenderID
		w.FunctionRequest = *FunctionRequest
		w.isTimeOut = false
		waiting <- w
		var client Client = e
		go waitingFunc(client, FunctionRequest, &FunctionRequest.BotConfig)
	}
	if isHotMusic.is {
		if isHotMusic.types != 1 {
			msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》没有在VtbMusic上找到结果。以下是VtbMusic的推荐:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
			cqfunction.CQSendMsg(FunctionRequest, msgMake)
			go do()
		} else {
			msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\nVTBMusic 当前已收录歌曲 " + strconv.Itoa(*isHotMusic.TotalQuantity) + "首。以下是VtbMusic的推荐:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
			cqfunction.CQSendMsg(FunctionRequest, msgMake)
			go do()
		}
	} else {
		switch FunctionRequest.MsgType {
		case "private":
			if lens == 1 {
				msgMake = e.getMusicDetailandCQCode(1)
			} else if lens <= 200 {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				go do()
			} else {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试,获取使用帮助请发送vtbhelp"
			}
			go cqfunction.CQSendPrivateMsg(FunctionRequest.SenderID, msgMake, &FunctionRequest.BotConfig)
			break
		case "group":
			if lens == 1 {
				msgMake = e.getMusicDetailandCQCode(1)
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgMake, &FunctionRequest.BotConfig)
			} else if lens <= 15 {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgMake, &FunctionRequest.BotConfig)
				go do()
			} else {
				if lens <= 40 {
					msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n请在原群聊发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
					msgtoGroup = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果。为防止打扰到他人，本消息采用私聊推送，请检查私信"
					go cqfunction.CQSendPrivateMsg(FunctionRequest.SenderID, msgMake, &FunctionRequest.BotConfig)
					do()
				} else {
					msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + MsgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试或私聊BOT获取完整列表,获取使用帮助请发送vtbhelp"
					msgtoGroup = msgMake
				}
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgtoGroup, &FunctionRequest.BotConfig)
			}
			break
		}

	}
}

func (e *VTBMusicClient) getMusicCode(Info *VTBMusicInfo) string {
	return fmt.Sprintf("[CQ:music,type=custom,url=https://vtbmusic.com/song?id=%s,audio=%s,title=%s,content=%s,image=%s]", Info.MusicID, Info.MusicURL, Info.MusicName, Info.MusicVocal, Info.Cover)
}

func (e *VTBMusicClient) getMusicDetailandCQCode(index int) string {
	return e.getMusicCode(e.getMusicDetail(index))
}
func (e *VTBMusicClient) musicListLen() int {
	return len(e.MusicList)
}
