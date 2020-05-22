package vtbmusic

import (
	"SMLKBOT/botstruct"
	"SMLKBOT/cqfunction"
	"encoding/json"
	"log"

	"github.com/tidwall/gjson"
)

//Reverse proxy for VtbmusicAPI
const vtbMusicSearchAPIAddr string = "https://api.smlk.org/Music_Manage/music_data/GetDataList"
const vtbMusicDetailAPIAddr string = "https://api.smlk.org/Music_Manage/music_data/GetTheData"
const vtbMusicCDNDetailAPIAddr string = "https://api.smlk.org/CDN_Manage/storage_data/GetDataList"

//GetVTBMusicList : Get VTBMusic Detail Info.
func GetVTBMusicList(musicname string) (VTBMusicList *botstruct.VTBMusicList) {
	ml := new(botstruct.VTBMusicList)
	s := make(map[string]string)
	s["condition"] = "name"
	s["keyword"] = musicname
	postjson := vtbPOSTJson{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result := cqfunction.WebPostJSONContent(vtbMusicSearchAPIAddr, string(p))
	ml.Total = gjson.GetBytes(result, "Total").Int()
	ml.Data = gjson.GetBytes(result, "Data").Array()
	return ml
}

//GetVTBMusicCDN : Get VTBMusic CDN Detail Info
func GetVTBMusicCDN(keyword string) (addr string) {
	//log.Println(keyword)
	s := make(map[string]string)
	s["condition"] = "name"
	s["keyword"] = keyword
	postjson := vtbCDNJson{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(p))
	result := cqfunction.WebPostJSONContent(vtbMusicCDNDetailAPIAddr, string(p))
	//log.Println(string(result))
	if gjson.GetBytes(result, "Data").IsArray() && len(gjson.GetBytes(result, "Data").Array()) != 0 {
		var addr string
		for _,r:= range gjson.GetBytes(result, "Data").Array() {
			if r.Get("name").String() == keyword{
				addr = r.Get("url").String()
				return addr
			}
		}
		
	}
	return ""
}

type vtbCDNJson struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
type vtbPOSTJson struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
