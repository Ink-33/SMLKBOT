package vtbmusic

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
	cpf = profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "nlp.tencentcloudapi.com"
	cpf.HttpProfile.ReqTimeout = 10
}

//Load : 加载配置
func Load() {
	credential = common.NewCredential(
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretId").String(),
		gjson.Get(*cqfunction.ConfigFile, "Tencent.secretKey").String(),
	)
}

//TenKeywordsExtraction :请求腾讯NLP API,传参格式为{"Text":"text"}
func TenKeywordsExtraction(params string) (result string) {
	client, err := nlp.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		log.Fatalln(err)
	}
	request := nlp.NewKeywordsExtractionRequest()
	err = request.FromJsonString(params)
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

//getNLPRequestString : 将Keyword格式化为Json
func getNLPRequestString(params string) (result string) {
	jsonstruct := &nlpRequest{
		KeyWord: params,
	}
	jsonbytes, err := json.Marshal(jsonstruct)
	if err != nil {
		log.Fatalln(err)
	}
	return string(jsonbytes)
}

type nlpRequest struct {
	KeyWord string `json:"Text"`
}
