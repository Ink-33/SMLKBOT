package txcloudutils

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"SMLKBOT/utils/smlkshell"
	"encoding/json"
	"log"
	"strings"
)

//TxCloudUtils : The massage handler of TxCloudUtils
func TxCloudUtils(FunctionRequest *botstruct.FunctionRequest) {
	msgArray := smlkshell.MsgSplit(FunctionRequest.Message)
	switch msgArray[0] {
	case "nlp":
		requeststr := strings.Join(msgArray[1:], " ")
		keywordJSON := TenKeywordsExtraction(requeststr)
		keywordStruct := new(KeywordsExtractionRespose)
		err := json.Unmarshal([]byte(keywordJSON), keywordStruct)
		if err != nil {
			log.Println(err)
			msgMake := "An unexpected error occurred while fetching data, please check console."
			cqfunction.CQSendMsg(FunctionRequest, msgMake)
			return
		}
		keywordArray := keywordStruct.Response.Keywords
		str := "[CQ:at,qq=" + FunctionRequest.SenderID + "]\n关键词提取结果如下:\n"
		tmpArray := make([]string, 0)
		for k, v := range keywordArray {
			tmpArray = append(tmpArray, string(k)+",分值:"+string(v.Score)+" - "+v.Word)
		}
		str = str + strings.Join(tmpArray, "\n")
		go cqfunction.CQSendMsg(FunctionRequest, str)
	}
}
