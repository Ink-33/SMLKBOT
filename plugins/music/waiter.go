package music

import (
	"log"
	"strconv"
	"time"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

var (
	waiting = make(chan *waitingChan, 15)
	counter int8
)

type newRequest struct {
	IsNewRequest    bool
	RequestSenderID string
}

type waitingChan struct {
	botstruct.FunctionRequest
	isTimeOut bool
	newRequest
}

func waitingFunc(client Client, fr *botstruct.FunctionRequest, botConfig *botstruct.BotConfig) {
	go func(fr *botstruct.FunctionRequest) {
		time.Sleep(45 * time.Second)
		wc := new(waitingChan)
		wc.FunctionRequest = *fr
		wc.isTimeOut = true
		waiting <- wc
	}(fr)
	for {
		c := <-waiting
		if c.IsNewRequest && c.RequestSenderID == fr.SenderID && c.HMACSHA1 != fr.HMACSHA1 {
			counter--
			return
		}
		if c.isTimeOut && c.HMACSHA1 == fr.HMACSHA1 {
			counter--
			return
		}
		if c.TimeStamp > fr.TimeStamp && isNumber(c.Message) {
			index, err := strconv.Atoi(c.Message)
			if err != nil {
				log.Fatalln(err)
			}
			if index <= client.musicListLen() && index > 0 {
				if c.SenderID == fr.SenderID && c.MsgType == fr.MsgType {
					log.SetPrefix("Music: ")
					log.Println("Known command:", fr.Message)
					switch c.MsgType {
					case "private":
						cqCodeMake := client.getMusicDetailandCQCode(index)
						counter--
						go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, botConfig)
					case "group":
						if c.GroupID == fr.GroupID {
							cqCodeMake := client.getMusicDetailandCQCode(index)
							counter--
							go cqfunction.CQSendGroupMsg(c.GroupID, cqCodeMake, botConfig)
						}
					}
					break
				}
			}
		}
	}
}
