package cqfunction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
)

// TimeOutError : A time out error
type TimeOutError struct {
	Addr string
}

func (e *TimeOutError) Error() string {
	return fmt.Sprint("Time out : " + e.Addr)
}

// ConfigFile : bot config file
var ConfigFile *string

func init() {
	ConfigFile = ReadConfig()
}

// CQSendGroupMsg : Send Group message by using CoolQ HTTPAPI
func CQSendGroupMsg(id, msg string, botConfig *botstruct.BotConfig) {
	msgstruct := &CQGroupMsg{ID: id, Message: msg}
	msgjson, err := json.Marshal(msgstruct)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = WebPostJSONContent(botConfig.HTTPAPIAddr+"/send_group_msg?access_token="+botConfig.HTTPAPIToken, string(msgjson))
	if err != nil {
		_, ok := err.(*TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}

// CQSendPrivateMsg : Send private message by using CoolQ HTTPAPI
func CQSendPrivateMsg(id, msg string, botConfig *botstruct.BotConfig) {
	msgstruct := &CQPrivateMsg{ID: id, Message: msg}
	msgjson, err := json.Marshal(msgstruct)
	if err != nil {
		log.Println(err.Error())
		return
	}
	_, err = WebPostJSONContent(botConfig.HTTPAPIAddr+"/send_private_msg?access_token="+botConfig.HTTPAPIToken, string(msgjson))
	if err != nil {
		_, ok := err.(*TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
}

// CQSendMsg : Send message
func CQSendMsg(functionRequest *botstruct.FunctionRequest, msg string) {
	switch functionRequest.MsgType {
	case "private":
		go CQSendPrivateMsg(functionRequest.SenderID, msg, &functionRequest.BotConfig)
	case "group":
		go CQSendGroupMsg(functionRequest.GroupID, msg, &functionRequest.BotConfig)
	}
}

// GetWebContent : Get web Content by using GET request.
func GetWebContent(addr string) (body []byte, err error) {
	content := make([]byte, 0)
	client := &http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}
	request, err := http.NewRequest("GET", addr, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4239.0 Safari/537.36")
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		_, ok := err.(net.Error)
		if ok {
			err1 := &TimeOutError{
				Addr: addr,
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

// WebPostJSONContent : Get web Content by using GET request.
func WebPostJSONContent(addr string, postbody string) (body []byte, err error) {
	content := make([]byte, 0)
	client := &http.Client{
		Transport: nil,
		Jar:       nil,
		Timeout:   10 * time.Second,
	}
	pb := strings.NewReader(postbody)
	request, err := http.NewRequest("POST", addr, pb)
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
				Addr: addr,
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

// ReadConfig : Read config file.
func ReadConfig() *string {
	file, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal("配置文件读取失败:", err)
	}
	result := string(file)
	return &result
}

// ReturnConfig : Return the config
func ReturnConfig() *string {
	return ConfigFile
}

// CQPrivateMsg : 私聊消息
type CQPrivateMsg struct {
	ID      string `json:"user_id"`
	Message string `json:"message"`
}

// CQGroupMsg : 群聊消息
type CQGroupMsg struct {
	ID      string `json:"group_id"`
	Message string `json:"message"`
}
