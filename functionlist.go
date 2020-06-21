package main

import (
	"SMLKBOT/biliau2card"
	"SMLKBOT/vtbmusic"
)

var functionList = make(map[string]functionFormat)

func functionLoad() {
	functionList["BiliAu2Card"] = biliau2card.Au2Card
	functionList["VTBMusic"] = vtbmusic.VTBMusic
}
