package txcloudutils

import (
	"SMLKBOT/utils/cqfunction"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"github.com/tidwall/gjson"
)

var credential *common.Credential
var nlpcpf *profile.ClientProfile

func init() {
	Load()
	nlpcpf = profile.NewClientProfile()
	nlpcpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	nlpcpf.HttpProfile.ReqTimeout = 10
}

//Load : 加载配置
func Load() {
	credential = common.NewCredential(
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretId").String(),
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretKey").String(),
	)
}

//TenKeywordsExtraction :请求腾讯TenKeywordsExtraction API,传入待提取文本与关键词数量上限
func TenKeywordsExtraction(text string, quantity uint64) (result string) {
	client, err := nlp.NewClient(credential, "ap-guangzhou", nlpcpf)
	if err != nil {
		log.Fatalln(err)
	}
	request := nlp.NewKeywordsExtractionRequest()
	request.Num = &quantity
	request.Text = &text
	response, err := client.KeywordsExtraction(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	result = response.ToJsonString()
	return
}

//KeywordsExtractionRequest : NLP文段摘要请求结构体
type KeywordsExtractionRequest struct {
	KeyWord string `json:"Text"`
}

//KeywordsExtractionRespose : NLP文段摘要返回值结构体
type KeywordsExtractionRespose struct {
	Response struct {
		Keywords []KeywordsExtractionKeywords `json:"Keywords"`
		Error    *struct {
			Code    string `json:"Code"`
			Message string `json:"Message"`
		} `json:"Error"`
		RequestID string `json:"RequestId"`
	} `json:"Response"`
}

//KeywordsExtractionKeywords : NLP文段摘要返回keywords字段结构体
type KeywordsExtractionKeywords struct {
	Score float64 `json:"Score"`
	Word  string  `json:"Word"`
}

//TenAutoSummarization :请求腾讯TenAutoSummarization API,传入待摘要文本与摘要的长度上限上限
func TenAutoSummarization(Text string, Length uint64) (summary string) {
	client, err := nlp.NewClient(credential, "ap-guangzhou", nlpcpf)
	if err != nil {
		log.Fatalln(err)
	}
	request := nlp.NewAutoSummarizationRequest()
	request.Text = &Text
	request.Length = &Length
	response, err := client.AutoSummarization(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
	return *response.Response.Summary
}
