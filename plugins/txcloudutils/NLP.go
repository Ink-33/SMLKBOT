package txcloudutils

import (
	"SMLKBOT/utils/cqfunction"
	"encoding/json"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	nlp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/nlp/v20190408"
	"github.com/tidwall/gjson"
)

var credential *common.Credential
var cpf *profile.ClientProfile

func init() {
	Load()
}

//Load : 加载配置
func Load() {
	credential = common.NewCredential(
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretId").String(),
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretKey").String(),
	)
}

//TenKeywordsExtraction :请求腾讯TenKeywordsExtraction API,直接传入待提取文本
func TenKeywordsExtraction(params string) (result string) {
	jsonstruct := &KeywordsExtractionRequest{
		KeyWord: params,
	}
	jsonbytes, err := json.Marshal(jsonstruct)
	if err != nil {
		log.Fatalln(err)
	}
	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	cpf.HttpProfile.ReqTimeout = 10
	client, err := nlp.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		log.Fatalln(err)
	}
	request := nlp.NewKeywordsExtractionRequest()
	err = request.FromJsonString(string(jsonbytes))
	if err != nil {
		log.Fatalln(err)
	}
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
	Score int    `json:"Score"`
	Word  string `json:"Word"`
}
