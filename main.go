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
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

var cqsecret string = gjson.Get(cqfunction.ReadConfig(), "HttpAPIPosSecret").String()

//Get audio number by regexp.
func getAu(msg string) (au string) {
	if strings.Contains(msg, "CQ:rich") {
		return ""
	}
	reg, err := regexp.Compile("(?i)au[0-9]+")
	if err != nil {
		log.Fatalln(err)
	}
	str := strings.Join(reg.FindAllString(msg, 1), "")
	return str
}

//Handle meassage and send music card.
func au2card(MsgInfo botstruct.MsgInfo) {
	au := getAu(MsgInfo.Message)
	if au != "" {
		log.Println("Created request for", au, "from:", MsgInfo.SenderID)
		Auinfo := biliau2card.GetAuInfo(au)

		if !Auinfo.AuStatus {
			msgMake := "AU" + Auinfo.AuNumber + Auinfo.AuMsg
			log.Println(msgMake)
			switch MsgInfo.MsgType {
			case "private":
				cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, msgMake)
				break
			case "group":
				cqfunction.CQSendGroupMsg(MsgInfo.GroupID, msgMake)
				break
			}
		} else {
			cqCodeMake := "[CQ:music,type=custom,url=" + Auinfo.AuJumpURL + ",audio=" + Auinfo.AuURL + ",title=" + Auinfo.AuTitle + ",content=" + Auinfo.AuDesp + ",image=" + Auinfo.AuCoverURL + "@180w_180h]"
			switch MsgInfo.MsgType {
			case "private":
				cqfunction.CQSendPrivateMsg(MsgInfo.SenderID, cqCodeMake)
				break
			case "group":
				cqfunction.CQSendGroupMsg(MsgInfo.GroupID, cqCodeMake)
				break
			}
		}

	} else {
		log.Println("Ingore message:", MsgInfo.Message, "from:", MsgInfo.SenderID)
	}
}

//handleMsg converts HTTP Post Body to MsgInfo Struct.
func handleMsg(raw []byte) (MsgInfo botstruct.MsgInfo) {
	MsgInfo.MsgType = gjson.GetBytes(raw, "message_type").String()
	MsgInfo.GroupID = gjson.GetBytes(raw, "group_id").String()
	MsgInfo.Message = gjson.GetBytes(raw, "message").String()
	MsgInfo.SenderID = gjson.GetBytes(raw, "user_id").String()
	return
}

//Handle request type before handling message.
func handleHTTP(w http.ResponseWriter, r *http.Request) {
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
		go au2card(handleMsg(body))
	}
}

func main() {
	config := cqfunction.ReadConfig()
	log.SetPrefix("SMLKBOT: ")
	path := gjson.Get(config, "BiliAu2Card.0.ListeningPath").String()
	port := gjson.Get(config, "BiliAu2Card.0.ListeningPort").String()
	log.Println("Powered by Ink33")
	log.Println("Start listening", path, port)

	http.HandleFunc(path, handleHTTP)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
