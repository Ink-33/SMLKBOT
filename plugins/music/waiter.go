package music

import (
	"SMLKBOT/data/botstruct"
	"SMLKBOT/utils/cqfunction"
	"log"
	"strconv"
	"time"
)

var waiting = make(chan *waitingChan, 15)
var counter int8 = 0

type newRequest struct {
	isNewRequest    bool
	RequestSenderID string
}
type waitingChan struct {
	botstruct.FunctionRequest
	isTimeOut bool
	newRequest
}

func waitingFunc(client Client, FunctionRequest *botstruct.FunctionRequest, BotConfig *botstruct.BotConfig) {
	go func(FunctionRequest *botstruct.FunctionRequest) {
		time.Sleep(45 * time.Second)
		wc := new(waitingChan)
		wc.FunctionRequest = *FunctionRequest
		wc.isTimeOut = true
		waiting <- wc
	}(FunctionRequest)
	for {
		c := <-waiting
		if c.isNewRequest && c.RequestSenderID == FunctionRequest.SenderID && c.HMACSHA1 != FunctionRequest.HMACSHA1 {
			counter--
			return
		}
		if c.isTimeOut && c.HMACSHA1 == FunctionRequest.HMACSHA1 {
			counter--
			return
		}
		if c.TimeStamp > FunctionRequest.TimeStamp && isNumber(c.Message) {
			index, err := strconv.Atoi(c.Message)
			if err != nil {
				log.Fatalln(err)
			}
			if index <= client.musicListLen() && index > 0 {
				if c.SenderID == FunctionRequest.SenderID && c.MsgType == FunctionRequest.MsgType {
					log.SetPrefix("Music: ")
					log.Println("Known command:", FunctionRequest.Message)
					switch c.MsgType {
					case "private":
						cqCodeMake := client.getMusicDetailandCQCode(index)
						counter--
						go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, BotConfig)
						break
					case "group":
						if c.GroupID == FunctionRequest.GroupID {
							cqCodeMake := client.getMusicDetailandCQCode(index)
							counter--
							go cqfunction.CQSendGroupMsg(c.GroupID, cqCodeMake, BotConfig)
							break
						}
					}
					break
				}
			}
		}
	}
}
