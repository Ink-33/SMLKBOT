package cqfunction

import (
	"SMLKBOT/botstruct"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//CQSendGroupMsg : Send Group message by using CoolQ HTTPAPI
func CQSendGroupMsg(id, msg string, BotConfig *botstruct.BotConfig) {
	//log.Println("cqFunctionGroup" + BotConfig.HTTPAPIAddr + "/send_group_msg?access_token=" + BotConfig.HTTPAPIToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
	GetWebContent(BotConfig.HTTPAPIAddr + "/send_group_msg?access_token=" + BotConfig.HTTPAPIToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
}

//CQSendPrivateMsg : Send private message by using CoolQ HTTPAPI
func CQSendPrivateMsg(id, msg string, BotConfig *botstruct.BotConfig) {
	//log.Println(BotConfig.HTTPAPIAddr + "/send_private_msg?access_token=" + BotConfig.HTTPAPIToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
	GetWebContent(BotConfig.HTTPAPIAddr + "/send_private_msg?access_token=" + BotConfig.HTTPAPIToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
}

//GetWebContent : Get web Content by using GET request.
func GetWebContent(url string) (body []byte) {
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

//WebPostJSONContent : Get web Content by using GET request.
func WebPostJSONContent(Addr string, postbody string) (body []byte) {
	client := &http.Client{}
	pb := strings.NewReader(postbody)
	request, err := http.NewRequest("POST", Addr, pb)
	request.Header.Set("Content-Type", "application/json")
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
