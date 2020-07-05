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
	postjson := vtbPOSTJSON{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(vtbMusicSearchAPIAddr, string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
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
	postjson := vtbCDNJSON{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}
	//log.Println(string(p))
	result, err := cqfunction.WebPostJSONContent(vtbMusicCDNDetailAPIAddr, string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
	//log.Println(string(result))
	if gjson.GetBytes(result, "Data").Value() != nil {
		var addr string
		for _, r := range gjson.GetBytes(result, "Data").Array() {
			if r.Get("name").String() == keyword {
				addr = r.Get("url").String()
				return addr
			}
		}
	}
	return ""
}

//GetVTBVocalList : Get VTB Detail Info.
func GetVTBVocalList(vocalname string) (VTBMusicList *botstruct.VTBMusicList) {
	ml := new(botstruct.VTBMusicList)
	s := make(map[string]string)
	s["condition"] = "vocal"
	s["keyword"] = vocalname
	postjson := vtbPOSTJSON{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(vtbMusicSearchAPIAddr, string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
	ml.Total = gjson.GetBytes(result, "Total").Int()
	ml.Data = gjson.GetBytes(result, "Data").Array()
	return ml
}

//GetVTBMusicDetail : Get music info by using music id.
func GetVTBMusicDetail(VTBMusicID string) (MusicInfo *botstruct.VTBMusicList) {
	ml := new(botstruct.VTBMusicList)
	postjson := vtbDetailJSON{
		MusicID: VTBMusicID,
	}
	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(vtbMusicDetailAPIAddr, string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
		} else {
			log.Fatalln(err)
		}
	}
	if gjson.GetBytes(result, "Data").Value() != nil {
		gr := make([]gjson.Result, 1)
		gr[0] = gjson.GetBytes(result, "Data|@this")
		ml.Data = gr
		ml.Total = 1
		return ml
	}
	ml.Total = 0
	return ml
}

type vtbCDNJSON struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
type vtbPOSTJSON struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
type vtbDetailJSON struct {
	MusicID string `json:"id"`
}
