package txcloudutils

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/data/helps"
	"SMLKBOT/utils/cqfunction"
	"SMLKBOT/utils/smlkshell"
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
)

//TxCloudUtils : The massage handler of TxCloudUtils
func TxCloudUtils(FunctionRequest *botstruct.FunctionRequest) {
	msgArray := smlkshell.MsgSplit(FunctionRequest.Message)
	if msgArray[0] != "txc" {
		return
	}
	switch msgArray[1] {
	case "--nlp":
		switch msgArray[2] {
		case "--keywordext":
			nlpKeywords(FunctionRequest, msgArray[3:])
			break
		default:
			cqfunction.CQSendMsg(FunctionRequest, "用法错误，请发送txc --help获取帮助")
			break
		}

	case "--help":
		sendHelp(FunctionRequest)
		break
	default:
		cqfunction.CQSendMsg(FunctionRequest, "用法错误，请发送txc --help获取帮助")
		break
	}

}
func sendHelp(FunctionRequest *botstruct.FunctionRequest) {
	cqfunction.CQSendMsg(FunctionRequest, help.TxCloudUtils)
}
func nlpKeywords(FunctionRequest *botstruct.FunctionRequest, msgArray []string) {
	quantity := isNumber(msgArray[0])
	if quantity == 0 {
		cqfunction.CQSendMsg(FunctionRequest, "用法错误，请发送txc --help获取帮助")
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
		cqfunction.CQSendMsg(FunctionRequest, msgMake)
		return
	}
	keywordArray := keywordStruct.Response.Keywords
	if len(keywordArray) != 0 {
		str := "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n关键词提取结果如下:\n"
		tmpArray := make([]string, 0)
		for k, v := range keywordArray {
			tmpArray = append(tmpArray, strconv.Itoa(k)+",分值:"+strconv.FormatFloat(v.Score, 'f', 8, 64)+" - "+v.Word)
		}
		go cqfunction.CQSendMsg(FunctionRequest, str+strings.Join(tmpArray, "\n"))
	} else {
		go cqfunction.CQSendMsg(FunctionRequest, "[CQ:at,qq="+FunctionRequest.SenderID+"]关键词提取结果为空。")
	}
}
func isNumber(str string) uint64 {
	var result uint64 = 0
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
