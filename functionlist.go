package main

import (
	"SMLKBOT/bilibili"
	"SMLKBOT/vtbmusic"
)

var functionList = make(map[string]functionFormat)

func functionLoad() {
	functionList["Bilibili"] = bilibili.Bilibili
	functionList["VTBMusic"] = vtbmusic.VTBMusic
}

func functionReload() {
	vtbmusic.Load()
}
