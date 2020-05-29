package main

import (
	"SMLKBOT/biliau2card"
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"SMLKBOT/vtbmusic"
	_ "crypto/hmac"
	_ "crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tidwall/gjson"
)

type function func(MsgInfo *botstruct.MsgInfo, BotConfig *botstruct.BotConfig)

var configfile string = cqfunction.ReadConfig()
var cqsecret string = gjson.Get(configfile, "HTTPAPIPostSecret").String()

func judgeandrun(name string, function function, MsgInfo *botstruct.MsgInfo) {
	var bc = new(botstruct.BotConfig)
	bc.HTTPAPIAddr = gjson.Get(configfile, "CoolQ.0.Api.HTTPAPIAddr").String()
	bc.HTTPAPIToken = gjson.Get(configfile, "CoolQ.0.Api.HTTPAPIToken").String()
	config := gjson.Get(configfile, "Feature.0").String()
	if gjson.Get(config, name).Bool() {
		go function(MsgInfo, bc)
	}
}

//MsgHandler converts HTTP Post Body to MsgInfo Struct.
func MsgHandler(raw []byte) (MsgInfo *botstruct.MsgInfo) {
	var mi = new(botstruct.MsgInfo)
	mi.TimeStamp = gjson.GetBytes(raw, "time").Int()
	mi.MsgType = gjson.GetBytes(raw, "message_type").String()
	mi.GroupID = gjson.GetBytes(raw, "group_id").String()
	mi.Message = gjson.GetBytes(raw, "message").String()
	mi.SenderID = gjson.GetBytes(raw, "user_id").String()

	return mi
}

//HTTPhandler : Handle request type before handling message.
func HTTPhandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method != "POST" {
		w.WriteHeader(400)
		fmt.Fprint(w, "Bad request.")
	} else {
		//signature := r.Header.Get("X-Signature")

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatalln(err)
		}
		//TODO: Add HMAC-SHA1 signature verification.
		/*
			mac1 := hmac.New(sha1.New, []byte(cqsecret))
			fmt.Println(string(body[:]))
			mac1.Write(body)
			fmt.Printf("%x\n",mac1.Sum(nil))
			fmt.Println(signature)
			fmt.Println(hmac.Equal(mac1.Sum(nil), []byte(signature)))
		*/
		var msgInfoTmp = MsgHandler(body)
		log.SetPrefix("SMLKBOT: ")
		log.Println("Received message:", msgInfoTmp.Message, "from:", msgInfoTmp.SenderID)
		go judgeandrun("BiliAu2Card", biliau2card.Au2Card, msgInfoTmp)
		go judgeandrun("VTBMusic", vtbmusic.VTBMusic, msgInfoTmp)
	}
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

func main() {
	log.SetPrefix("SMLKBOT: ")
	closeSignalHandler()
	path := gjson.Get(configfile, "CoolQ.0.HTTPServer.ListeningPath").String()
	port := gjson.Get(configfile, "CoolQ.0.HTTPServer.ListeningPort").String()

	log.Println("Powered by Ink33")
	log.Println("Start listening", path, port)

	http.HandleFunc(path, HTTPhandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
