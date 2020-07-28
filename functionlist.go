package main

import (
	"SMLKBOT/plugins/bilibili"
	"SMLKBOT/plugins/vtbmusic"
)

var functionList = make(map[string]functionFormat)

func functionLoad() {
	functionList["Bilibili"] = bilibili.Bilibili
	functionList["VTBMusic"] = vtbmusic.VTBMusic
}

func functionReload() {
	vtbmusic.Load()
}
