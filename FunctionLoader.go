package main

import (
	"SMLKBOT/plugins/bilibili"
	"SMLKBOT/plugins/txcloudutils"
	"SMLKBOT/plugins/music"
)

var functionList = make(map[string]functionFormat)

func init() {
	functionList["Bilibili"] = bilibili.Bilibili
	functionList["Music"] = music.MsgHandler
	functionList["TxCloudUtils"] = txcloudutils.TxCloudUtils
}

