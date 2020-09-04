package main

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"SMLKBOT/utils/smlkshell"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tidwall/gjson"
)

type functionFormat func(*botstruct.FunctionRequest)

func judgeandrun(name string, functionFormat functionFormat, FunctionRequest *botstruct.FunctionRequest) {
	config := gjson.Get(*cqfunction.ConfigFile, "Feature.0").String()
	if gjson.Get(config, name).Bool() {
		go functionFormat(FunctionRequest)
	}
}

//MsgHandler converts HTTP Post Body to MsgInfo Struct.
func MsgHandler(hmacsha1 string, raw []byte) (MsgInfo *botstruct.MsgInfo) {
	var mi = new(botstruct.MsgInfo)
	mi.TimeStamp = gjson.GetBytes(raw, "time").Int()
	mi.MsgType = gjson.GetBytes(raw, "message_type").String()
	mi.GroupID = gjson.GetBytes(raw, "group_id").String()
	mi.Message = gjson.GetBytes(raw, "message").String()
	mi.SenderID = gjson.GetBytes(raw, "user_id").String()
	mi.Role = gjson.GetBytes(raw, "sender.role").String()
	mi.HMACSHA1 = hmacsha1
	return mi
}

//HTTPhandler : Handle request type before handling message.
func HTTPhandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "POST" {
		w.WriteHeader(400)
		fmt.Fprint(w, "<body><img src=\"https://api.smlk.org/mirror/ink33/what.png\" style=\"vertical-align: top\" alt=\"Bad request.\"/><ln><p>Bad request.</p></body>")
	} else {
		rid := r.Header.Get("X-Self-ID")
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		hmacSH1 := hmac.New(sha1.New, []byte(gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+rid+".HTTPAPIPostSecret").String()))
		hmacSH1.Reset()
		hmacSH1.Write(body)
		var signature string = strings.Replace(r.Header.Get("X-Signature"), "sha1=", "", 1)
		var hmacresult string = fmt.Sprintf("%x", hmacSH1.Sum(nil))
		if signature == "" {
			w.WriteHeader(401)
			fmt.Fprint(w, "Unauthorized.")
		} else if signature != hmacresult {
			w.WriteHeader(401)
			fmt.Fprint(w, "Unauthorized.")
		} else {
			if gjson.GetBytes(body, "meta_event_type").String() != "heartbeat" {
				var msgInfoTmp = MsgHandler(signature, body)
				msgInfoTmp.RobotID = rid
				var bc = new(botstruct.BotConfig)
				bc.HTTPAPIAddr = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+msgInfoTmp.RobotID+".HTTPAPIAddr").String()
				bc.HTTPAPIToken = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+msgInfoTmp.RobotID+".HTTPAPIToken").String()
				bc.MasterID = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Master").Array()
				log.SetPrefix("SMLKBOT: ")
				fr := &botstruct.FunctionRequest{MsgInfo: *msgInfoTmp, BotConfig: *bc}
				go log.Println("RobotID:", rid, "Received message:", msgInfoTmp.Message, "from:", "User:", msgInfoTmp.SenderID, "Group:", msgInfoTmp.GroupID, "Role:", smlkshell.RoleHandler(fr).RoleName)
				smlkshell.SmlkShell(fr)
				for k, v := range functionList {
					go judgeandrun(k, v, fr)
				}
			}
		}
	}
}

func scfHandler(event apigwEvent) (result *scfReturn, err error) {
	if event.Headers.BotID == "" {
		return newSCFReturn(406, "Not Acceptable."), nil
	}
	hmacSH1 := hmac.New(sha1.New, []byte(gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+event.Headers.BotID+".HTTPAPIPostSecret").String()))
	hmacSH1.Reset()
	hmacSH1.Write([]byte(event.Body))
	var signature string = strings.Replace(event.Headers.Signature, "sha1=", "", 1)
	var hmacresult string = fmt.Sprintf("%x", hmacSH1.Sum(nil))
	if signature == "" || signature != hmacresult {
		return newSCFReturn(401, "Unauthorized."), nil
	}
	if gjson.Get(event.Body, "meta_event_type").String() != "heartbeat" {
		var msgInfoTmp = MsgHandler(signature, []byte(event.Body))
		msgInfoTmp.RobotID = event.Headers.BotID
		var bc = new(botstruct.BotConfig)
		bc.HTTPAPIAddr = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+msgInfoTmp.RobotID+".HTTPAPIAddr").String()
		bc.HTTPAPIToken = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Api."+msgInfoTmp.RobotID+".HTTPAPIToken").String()
		bc.MasterID = gjson.Get(*cqfunction.ConfigFile, "CoolQ.Master").Array()
		log.SetPrefix("SMLKBOT: ")
		fr := &botstruct.FunctionRequest{MsgInfo: *msgInfoTmp, BotConfig: *bc}
		go log.Println("RobotID:", event.Headers.BotID, "Received message:", msgInfoTmp.Message, "from:", "User:", msgInfoTmp.SenderID, "Group:", msgInfoTmp.GroupID, "Role:", smlkshell.RoleHandler(fr).RoleName)
		smlkshell.SmlkShell(fr)
		for k, v := range functionList {
			go judgeandrun(k, v, fr)
		}

	}
	return newSCFReturn(200, "Success."), nil
}

func newSCFReturn(status int, msg string) *scfReturn {
	r := &scfReturn{
		Status: status,
		Msg:    msg,
	}
	if status != 200 {
		log.Println(*r)
	}
	return r
}
func closeSignalHandler() {
	channel := make(chan os.Signal, 2)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-channel
		log.Println("Program stop.")
		os.Exit(0)
	}()
}

func newHTTPServer() {
	closeSignalHandler()
	path := gjson.Get(*cqfunction.ConfigFile, "CoolQ.HTTPServer.ListeningPath").String()
	port := gjson.Get(*cqfunction.ConfigFile, "CoolQ.HTTPServer.ListeningPort").String()

	log.Println("Powered by Ink33")
	log.Println("Start listening", path, port)

	http.HandleFunc(path, HTTPhandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalln("ListenAndServe", err)
	}
}

func main() {
	log.SetPrefix("SMLKBOT: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	runtime.GOMAXPROCS(runtime.NumCPU())

	switch smlkshell.IsSCF {
	case "TencentSCF":
		cloudfunction.Start(scfHandler)
	default:
		newHTTPServer()
	}

}

type apigwEvent struct {
	Headers headers `json:"headers"`
	Body    string  `json:"body"`
}

type headers struct {
	BotID     string `json:"X-Self-ID"`
	Signature string `json:"X-Signature"`
}

type scfReturn struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}
