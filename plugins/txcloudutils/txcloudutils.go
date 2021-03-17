package txcloudutils

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	help "github.com/Ink-33/SMLKBOT/data/helps"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
	"github.com/Ink-33/SMLKBOT/utils/smlkshell"
)

// TxCloudUtils : The massage handler of TxCloudUtils
func TxCloudUtils(fr *botstruct.FunctionRequest) {
	msgArray := smlkshell.MsgSplit(fr.Message)
	if msgArray[0] != "txc" {
		return
	}
	switch msgArray[1] {
	case "--nlp":
		switch msgArray[2] {
		case "--keywordext":
			nlpKeywords(fr, msgArray[3:])
		case "--summarization":
			nlpSummarization(fr, msgArray[3:])
		default:
			cqfunction.CQSendMsg(fr, "用法错误，请发送txc --help获取帮助")
		}

	case "--help":
		sendHelp(fr)
	default:
		cqfunction.CQSendMsg(fr, "用法错误，请发送txc --help获取帮助")
	}
}

func sendHelp(fr *botstruct.FunctionRequest) {
	cqfunction.CQSendMsg(fr, help.TxCloudUtils)
}

func nlpKeywords(fr *botstruct.FunctionRequest, msgArray []string) {
	quantity := isNumber(msgArray[0])
	if quantity == 0 {
		cqfunction.CQSendMsg(fr, "用法错误，请发送txc --help获取帮助")
		return
	}
	log.Println("Known command: TenKeyWordExtraction")
	requeststr := strings.Join(msgArray[1:], " ")
	keywordJSON := TenKeywordsExtraction(requeststr, quantity)
	keywordStruct := new(KeywordsExtractionRespose)
	err := json.Unmarshal([]byte(keywordJSON), keywordStruct)
	if err != nil {
		log.Println(err)
		msgMake := "An unexpected error occurred while fetching data, please check console."
		cqfunction.CQSendMsg(fr, msgMake)
		return
	}
	keywordArray := keywordStruct.Response.Keywords
	if len(keywordArray) != 0 {
		str := "[CQ:at,qq=" + fr.SenderID + "]\n关键词提取结果如下:\n"
		tmpArray := make([]string, 0)
		for k, v := range keywordArray {
			tmpArray = append(tmpArray, strconv.Itoa(k)+",分值:"+strconv.FormatFloat(v.Score, 'f', 8, 64)+" - "+v.Word)
		}
		go cqfunction.CQSendMsg(fr, str+strings.Join(tmpArray, "\n"))
	} else {
		go cqfunction.CQSendMsg(fr, "[CQ:at,qq="+fr.SenderID+"]关键词提取结果为空。")
	}
}

func nlpSummarization(fr *botstruct.FunctionRequest, msgArray []string) {
	length := isNumber(msgArray[0])
	if length == 0 {
		cqfunction.CQSendMsg(fr, "用法错误，请发送txc --help获取帮助")
		return
	}
	log.Println("Known command: TenAutoSummarization")
	text := strings.Join(msgArray[1:], " ")
	summarization := TenAutoSummarization(text, length)
	if len(summarization) == 0 {
		cqfunction.CQSendMsg(fr, "摘要失败。")
		return
	}
	go cqfunction.CQSendMsg(fr, "[CQ:at,qq="+fr.SenderID+"]摘要结果如下:\n"+summarization)
}

func isNumber(str string) uint64 {
	var result uint64
	reg := regexp.MustCompile("^[0-9]+$")
	tmp := strings.Join(reg.FindAllString(str, 1), "")
	if tmp != "" {
		result, err := strconv.ParseUint(tmp, 10, 64)
		if err != nil {
			log.Println(err.Error())
			return 0
		}
		return result
	}
	return result
}
