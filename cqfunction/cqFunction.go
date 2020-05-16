package cqfunction

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

//CQSendGroupMsg : Send Group message by using CoolQ HttpAPI
func CQSendGroupMsg(id, msg string) {
	conf := ReadConfig()
	cqAddr := gjson.Get(conf, "CoolQ.0.Api.HttpAPIAddr").String()
	cqToken := gjson.Get(conf, "CoolQ.0.Api.HttpAPIToken").String()
	GetWbeContent(cqAddr + "/send_group_msg?access_token=" + cqToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
}

//CQSendPrivateMsg : Send private message by using CoolQ HttpAPI
func CQSendPrivateMsg(id, msg string) {
	conf := ReadConfig()
	cqAddr := gjson.Get(conf, "CoolQ.0.Api.HttpAPIAddr").String()
	cqToken := gjson.Get(conf, "CoolQ.0.Api.HttpAPIToken").String()
	GetWbeContent(cqAddr + "/send_private_msg?access_token=" + cqToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
}

//GetWbeContent : Get web Content by using GET request.
func GetWbeContent(url string) (body []byte) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4117.2 Safari/537.36")
	if err != nil {
		log.Fatalln(err)
	}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()
	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return content
}

//ReadConfig : Read config file.
func ReadConfig() string {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	result := string(file)
	return result
}
