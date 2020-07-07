package cqfunction

import (
	"SMLKBOT/botstruct"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//TimeOutError : A time out error
type TimeOutError struct {
	Addr string
}

func (e *TimeOutError) Error() string {

	return fmt.Sprint("Time out : " + e.Addr)
}

//ConfigFile : bot config file
var ConfigFile *string

func init() {
	ConfigFile = ReadConfig()
}

//CQSendGroupMsg : Send Group message by using CoolQ HTTPAPI
func CQSendGroupMsg(id, msg string, BotConfig *botstruct.BotConfig) {
	//log.Println("cqFunctionGroup" + BotConfig.HTTPAPIAddr + "/send_group_msg?access_token=" + BotConfig.HTTPAPIToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
	_, err := GetWebContent(BotConfig.HTTPAPIAddr + "/send_group_msg?access_token=" + BotConfig.HTTPAPIToken + "&group_id=" + id + "&message=" + url.QueryEscape(msg))
	if err != nil {
		_, ok := err.(*TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}

//CQSendPrivateMsg : Send private message by using CoolQ HTTPAPI
func CQSendPrivateMsg(id, msg string, BotConfig *botstruct.BotConfig) {
	//log.Println(BotConfig.HTTPAPIAddr + "/send_private_msg?access_token=" + BotConfig.HTTPAPIToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
	_, err := GetWebContent(BotConfig.HTTPAPIAddr + "/send_private_msg?access_token=" + BotConfig.HTTPAPIToken + "&user_id=" + id + "&message=" + url.QueryEscape(msg))
	if err != nil {
		_, ok := err.(*TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}

//CQSendMsg : Send message
func CQSendMsg(MsgInfo *botstruct.MsgInfo, msg string, BotConfig *botstruct.BotConfig) {
	switch MsgInfo.MsgType {
	case "private":
		go CQSendPrivateMsg(MsgInfo.SenderID, msg, BotConfig)
		break
	case "group":
		go CQSendGroupMsg(MsgInfo.GroupID, msg, BotConfig)
		break
	}
}

//GetWebContent : Get web Content by using GET request.
func GetWebContent(Addr string) (body []byte, err error) {
	var content = make([]byte, 0)
	client := &http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}
	request, err := http.NewRequest("GET", Addr, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4117.2 Safari/537.36")
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		_, ok := err.(net.Error)
		if ok {
			err1 := &TimeOutError{
				Addr: Addr,
			}
			return content, err1
		}
	}
	defer response.Body.Close()
	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//WebPostJSONContent : Get web Content by using GET request.
func WebPostJSONContent(Addr string, postbody string) (body []byte, err error) {
	var content = make([]byte, 0)
	client := &http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}
	pb := strings.NewReader(postbody)
	request, err := http.NewRequest("POST", Addr, pb)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4117.2 Safari/537.36")
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		_, ok := err.(net.Error)
		if ok {
			err1 := &TimeOutError{
				Addr: Addr,
			}
			return content, err1
		}
	}
	defer response.Body.Close()
	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return content, nil
}

//ReadConfig : Read config file.
func ReadConfig() *string {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	result := string(file)
	return &result
}

//ReturnConfig : Return the config
func ReturnConfig() *string {
	return ConfigFile
}
