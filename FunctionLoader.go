package main

import (
	"github.com/Ink-33/SMLKBOT/plugins/bilibili"
	"github.com/Ink-33/SMLKBOT/plugins/music"
	"github.com/Ink-33/SMLKBOT/plugins/txcloudutils"
)

var functionList = make(map[string]functionFormat)

func init() {
	functionList["Bilibili"] = bilibili.Bilibili
	functionList["Music"] = music.MsgHandler
	functionList["TxCloudUtils"] = txcloudutils.TxCloudUtils
}
