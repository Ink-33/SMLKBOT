package main

import (
	"SMLKBOT/biliau2card"
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
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

type function func(MsgInfo *botstruct.MsgInfo)

var cqsecret string = gjson.Get(cqfunction.ReadConfig(), "HttpAPIPosSecret").String()

func judgeandrun(name string, function function, MsgInfo *botstruct.MsgInfo) {
	config := gjson.Get(cqfunction.ReadConfig(), "Feature.0").String()

	if gjson.Get(config, name).Bool() {
		function(MsgInfo)
	} else {
		log.Println("Ingore message:", MsgInfo.Message, "from:", MsgInfo.SenderID)
	}
}

//MsgHandler converts HTTP Post Body to MsgInfo Struct.
func MsgHandler(raw []byte) (MsgInfo *botstruct.MsgInfo) {
	var mi = new(botstruct.MsgInfo)

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
		go judgeandrun("BiliAu2Card", biliau2card.Au2Card, msgInfoTmp)
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
	config := cqfunction.ReadConfig()
	log.SetPrefix("SMLKBOT: ")
	closeSignalHandler()

	path := gjson.Get(config, "CoolQ.0.HTTPServer.ListeningPath").String()
	port := gjson.Get(config, "CoolQ.0.HTTPServer.ListeningPort").String()

	log.Println("Powered by Ink33")
	log.Println("Start listening", path, port)

	http.HandleFunc(path, HTTPhandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
