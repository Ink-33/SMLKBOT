package vtbmusic

import (
	"strings"
	"SMLKBOT/cqfunction"
	"log"
)
//Reverse proxy for VtbmusicAPI
const vtbMusicSearchAPIAddr string = "https://api.smlk.org/Music_Manage/music_data/GetDataList"
const vtbMusicDetailAPIAddr string = "https://api.smlk.org/Music_Manage/music_data/GetTheData"
//T is a test Function
func T() {
	var tmp string = "{'search':{'condition':'name','keyword':'だんご大家族'},'pageIndex':1,'pageRows':9999}"
	log.Println(string(cqfunction.WebPostJSONContent(vtbMusicSearchAPIAddr, strings.Replace(tmp,"'","\"",999))))
}
