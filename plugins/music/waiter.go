package music

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/Ink-33/SMLKBOT/data/botstruct"
	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

var eventBus = &EventBus{
	subscribers: map[string]*EventBusSubcriber{},
	lock:        sync.RWMutex{},
}

func (e *EventBus) subscribe(client Client, fr *botstruct.FunctionRequest, botConfig *botstruct.BotConfig) {
	do := func(c *botstruct.FunctionRequest) {
		index, err := strconv.Atoi(c.Message)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetPrefix("Music: ")
		log.Println("Known command: Choose Music", c.Message)
		switch c.MsgType {
		case "private":
			cqCodeMake := client.getMusicDetailandCQCode(index)
			go cqfunction.CQSendPrivateMsg(c.SenderID, cqCodeMake, botConfig)
		case "group":
			if c.GroupID == fr.GroupID {
				cqCodeMake := client.getMusicDetailandCQCode(index)
				go cqfunction.CQSendGroupMsg(c.GroupID, cqCodeMake, botConfig)
			}
		}
		e.subscribers[c.SenderID].done <- struct{}{}
		e.lock.Lock()
		delete(e.subscribers, c.SenderID)
		e.lock.Unlock()
	}

	subscriber := &EventBusSubcriber{
		run:  do,
		done: make(chan struct{}, 1),
	}
	go func() {
		select {
		case <-time.After(45 * time.Second):
			e.lock.Lock()
			delete(e.subscribers, fr.SenderID)
			e.lock.Unlock()
		case <-subscriber.done:
			return
		}
	}()
	e.lock.Lock()
	e.subscribers[fr.SenderID] = subscriber
	e.lock.Unlock()
}
