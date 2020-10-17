package vtbmusic

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/data/helps"
	txc "SMLKBOT/plugins/txcloudutils"
	"SMLKBOT/utils/cqfunction"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//vtbmusic MsgType
type vtbMsgType struct {
	//Search Keyword
	content string
	ctype   int8
}

type isHotMusic struct {
	is bool
	/*
		0 -> 未找到歌曲
		1 -> 仅推荐
	*/
	types         int8
	TotalQuantity *int
}

var waiting = make(chan *waitingChan, 15)
var counter int8 = 0

type newRequest struct {
	isNewRequest    bool
	RequestSenderID string
}
type waitingChan struct {
	botstruct.FunctionRequest
	isTimeOut bool
	newRequest
}

//VTBMusic : The main function of VTBMusic
func VTBMusic(FunctionRequest *botstruct.FunctionRequest) {
	mt := msgHandler(FunctionRequest.Message)
	ctype := mt.ctype
	var isHot *isHotMusic
	switch ctype {
	case 0:
		break
	case 1:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command:", mt.content)
		go cqfunction.CQSendMsg(FunctionRequest, "Searching...")
		keywordjson := txc.TenKeywordsExtraction(mt.content)
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
			list1 := GetVTBMusicList(mt.content, "MusicName")
			list2 := GetVTBMusicList(mt.content, "VtbName")
			if list1.Total == -1 || list2.Total == -1 {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				cqfunction.CQSendMsg(FunctionRequest, msgMake)
				return
			}
			ListMsg, ListArray := listToMsg(list1, list2)
			if len(ListArray) == 0 {
				list := GetHotMusicList()
				ListMsg, ListArray := listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			}
		} else {
			nlpMsg, nlpArray := nlpListToMsg(keywordArray)
			if nlpArray == nil {
				msgMake := "An unexpected error occurred while fetching data, please check console."
				log.Println(msgMake)
				cqfunction.CQSendMsg(FunctionRequest, msgMake)
				return
			}
			if len(nlpArray) == 0 {
				list := GetHotMusicList()
				if list.Total == -1 {
					msgMake := "An unexpected error occurred while fetching data, please check console."
					log.Println(msgMake)
					cqfunction.CQSendMsg(FunctionRequest, msgMake)
					return
				}
				ListMsg, ListArray := listToMsg(list)
				isHot = &isHotMusic{true, 0, nil}
				sendMsg(FunctionRequest, ListMsg, ListArray, mt, isHot)
			} else {
				isHot = &isHotMusic{false, 0, nil}
				sendMsg(FunctionRequest, nlpMsg, nlpArray, mt, isHot)
			}
		}
		break
	case 2:
		log.SetPrefix("VTBMusic: ")
		log.Println("Known command: Get hot music.")
		list := GetHotMusicList()
		ListMsg, ListArray := listToMsg(list)
		var msgMake string
		if list.Total == -1 {
			msgMake = "An unexpected error occurred while fetching data, please check console."
			log.Println(msgMake)
		} else {
			isHot = &isHotMusic{true, 1, &list.Total}
			sendMsg(FunctionRequest, ListMsg, ListArray, nil, isHot)
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
		list := GetVTBMusicDetail(mt.content)
		var msgMake string
		if list.Total == -1 {
			msgMake = "An unexpected error occurred while fetching data, please check console."
			log.Println(msgMake)
		} else if list.Total == 0 {
			msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\nid:" + mt.content + "没有在VtbMusic上找到结果。获取使用帮助请发送vtbhelp"
		} else {
			info := getMusicDetail(list.Data, 1)
			msgMake = getMusicCode(info)
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

func msgHandler(msg string) (MsgType *vtbMsgType) {
	mt := new(vtbMsgType)
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

func listToMsg(list ...*MusicList) (ListMsg *string, ListArray []GetMusicListData) {
	var q = make([]string, 0)
	var listReturn = make([]GetMusicListData, 0)
	for _, v := range list {
		for i, r := range v.Data {
			listReturn = append(listReturn, r)
			t := strconv.Itoa(i+1) + "," + r.VocalName + "-" + r.OriginName
			q = append(q, t)
		}
	}
	msg := strings.Join(q, "\n")
	return &msg, listReturn
}

func waitingFunc(list []GetMusicListData, FunctionRequest *botstruct.FunctionRequest, BotConfig *botstruct.BotConfig) {
	go func(FunctionRequest *botstruct.FunctionRequest) {
		time.Sleep(45 * time.Second)
		wc := new(waitingChan)
		wc.FunctionRequest = *FunctionRequest
		wc.isTimeOut = true
		waiting <- wc
	}(FunctionRequest)
	for {
		c := <-waiting
		if c.isNewRequest && c.RequestSenderID == FunctionRequest.SenderID && c.HMACSHA1 != FunctionRequest.HMACSHA1 {
			counter--
			runtime.Goexit()
		}
		if c.isTimeOut && c.HMACSHA1 == FunctionRequest.HMACSHA1 {
			counter--
			runtime.Goexit()
		} else if c.TimeStamp > FunctionRequest.TimeStamp && isNumber(c.Message) {
			index, err := strconv.Atoi(c.Message)
			if err != nil {
				log.Fatalln(err)
			}
			if index <= len(list) && index > 0 {
				if c.SenderID == FunctionRequest.SenderID && c.MsgType == FunctionRequest.MsgType {
					switch c.MsgType {
					case "private":
						info := getMusicDetail(list, index)
						cqCodeMake := getMusicCode(info)
						counter--
						go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, BotConfig)
						break
					case "group":
						if c.GroupID == FunctionRequest.GroupID {
							info := getMusicDetail(list, index)
							cqCodeMake := getMusicCode(info)
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

func getMusicDetail(list []GetMusicListData, index int) (info *MusicInfo) {
	i := new(MusicInfo)
	i.MusicID = list[index-1].ID
	i.MusicVocal = strings.ReplaceAll(list[index-1].VocalName, ",", "、")
	i.MusicName = list[index-1].OriginName
	cdn := GetVTBMusicCDN("")
	if strings.Contains(list[index-1].CDN, ":") {
		reg := regexp.MustCompile("(\\d+):(\\d+):(\\d+)")
		match := reg.FindStringSubmatch(list[index-1].CDN)
		i.Cover = cdn.match(match[1]) + list[index-1].CoverImg
		i.MusicURL = cdn.match(match[2]) + list[index-1].Music
	} else {
		i.MusicCDN = cdn.match(list[index-1].CDN)
		i.Cover = i.MusicCDN + list[index-1].CoverImg
		i.MusicURL = i.MusicCDN + list[index-1].Music
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

func nlpListToMsg(keywordArray []txc.KeywordsExtractionKeywords) (NLPMsg *string, NLPArray []GetMusicListData) {
	list1 := GetVTBMusicList(keywordArray[0].Word, "MusicName")
	list2 := GetVTBMusicList(keywordArray[0].Word, "VtbName")

	if list1.Total == -1 || list2.Total == -1 {
		msgMake := "An unexpected error occurred while fetching data, please check console."
		log.Printf("Fail while getting music lists, list1:%d, list2:%d", list1.Total, list2.Total)
		return &msgMake, nil
	}

	var nlpArray = make([]GetMusicListData, 0)
	var nlpMsgArray = make([]string, 0)

	if list1.Total+list2.Total == 0 {
		msgMake := "nope"
		return &msgMake, nlpArray
	}

	_, ListArray := listToMsg(list1, list2)

	for k1, v1 := range keywordArray {
		if len(keywordArray) == 1 {
			nlpArray = ListArray
			break
		}
		switch k1 {
		case 1:
			reg := regexp.MustCompile("(?i)(" + v1.Word + ")")
			for _, v2 := range ListArray {
				if reg.MatchString(v2.VocalName) || reg.MatchString(v2.OriginName) {
					nlpArray = append(nlpArray, v2)
				}
			}
		default:
			reg := regexp.MustCompile("(?i)(" + v1.Word + ")")
			for _, v2 := range nlpArray {
				if reg.MatchString(v2.VocalName) || reg.MatchString(v2.OriginName) {
					nlpArray = append(nlpArray, v2)
				}
			}
		}

	}
	for k, v := range nlpArray {
		tmp := strconv.Itoa(k+1) + "," + v.VocalName + "-" + v.OriginName
		nlpMsgArray = append(nlpMsgArray, tmp)
	}
	msg := strings.Join(nlpMsgArray, "\n")
	return &msg, nlpArray
}

func sendMsg(FunctionRequest *botstruct.FunctionRequest, ListMsg *string, ListArray []GetMusicListData, vtbMsgType *vtbMsgType, isHotMusic *isHotMusic) {
	var msgMake string
	var msgtoGroup string
	lens := len(ListArray)
	do := func() {
		counter++
		w := new(waitingChan)
		w.isNewRequest = true
		w.RequestSenderID = FunctionRequest.SenderID
		w.FunctionRequest = *FunctionRequest
		w.isTimeOut = false
		waiting <- w
		go waitingFunc(ListArray, FunctionRequest, &FunctionRequest.BotConfig)
	}
	if isHotMusic.is {
		if isHotMusic.types != 1 {
			msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》没有在VtbMusic上找到结果。以下是VtbMusic的推荐:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
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
				info := getMusicDetail(ListArray, 1)
				msgMake = getMusicCode(info)
			} else if lens <= 200 {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				go do()
			} else {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试,获取使用帮助请发送vtbhelp"
			}
			go cqfunction.CQSendPrivateMsg(FunctionRequest.SenderID, msgMake, &FunctionRequest.BotConfig)
			break
		case "group":
			if lens == 1 {
				info := getMusicDetail(ListArray, 1)
				msgMake = getMusicCode(info)
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgMake, &FunctionRequest.BotConfig)
			} else if lens <= 15 {
				msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgMake, &FunctionRequest.BotConfig)
				go do()
			} else {
				if lens <= 40 {
					msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果:\n" + *ListMsg + "\n━━━━━━━━━━━━━━\n请在原群聊发送歌曲对应序号即可播放,获取使用帮助请发送vtbhelp"
					msgtoGroup = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到" + strconv.Itoa(lens) + "个结果。为防止打扰到他人，本消息采用私聊推送，请检查私信"
					go cqfunction.CQSendPrivateMsg(FunctionRequest.SenderID, msgMake, &FunctionRequest.BotConfig)
					do()
				} else {
					msgMake = "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n《" + vtbMsgType.content + "》共找到多达" + strconv.Itoa(lens) + "个结果,建议您更换关键词重试或私聊BOT获取完整列表,获取使用帮助请发送vtbhelp"
					msgtoGroup = msgMake
				}
				go cqfunction.CQSendGroupMsg(FunctionRequest.GroupID, msgtoGroup, &FunctionRequest.BotConfig)
			}
			break
		}

	}
}

func getMusicCode(Info *MusicInfo) string {
	return fmt.Sprintf("[CQ:music,type=custom,url=https://vtbmusic.com/song?id=%s,audio=%s,title=%s,content=%s,image=%s]", Info.MusicID, Info.MusicURL, Info.MusicName, Info.MusicVocal, Info.Cover)
}
