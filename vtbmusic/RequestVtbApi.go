package vtbmusic

import (
	"SMLKBOT/cqfunction"
	"encoding/json"
	"log"
	"runtime"

	"github.com/tidwall/gjson"
)

//Reverse proxy for VtbmusicAPI: https://api.vtbmusic.com:60006/v1/
const vtbMusicAPIProxy string = "https://api.smlk.org/v1/"

var vtbMusicAPIMap map[string]string = make(map[string]string)

func init() {
	vtbMusicAPIMap["MusicList"] = "GetMusicList"
	vtbMusicAPIMap["MusicData"] = "GetMusicData"
	vtbMusicAPIMap["VtbsList"] = "GetVtbsList"
	vtbMusicAPIMap["VtbsData"] = "GetVtbsData"
	vtbMusicAPIMap["CDN"] = "GetCDNList"
	vtbMusicAPIMap["HotList"] = "GetHotMusicList"
}
func getAPIAddr(APIName string) string {
	return vtbMusicAPIProxy + vtbMusicAPIMap[APIName]
}

//GetVTBMusicList : Get VTBMusic Detail Info.
//	Method: MusicName,VtbName
func GetVTBMusicList(Keyword string, Method string) (VTBMusicList *MusicList) {
	ml := new(MusicList)
	switch Method {
	case "VtbName":
		vl := GetVTBsList(Keyword)
		vidraw := make([]string, 0)
		for _, v := range vl.Data {
			vidraw = append(vidraw, v.Get("Id").String())
		}
		var vid []string
		if len(vidraw) > 3 {
			vid = vidraw[0:2]
		} else {
			vid = vidraw
		}
		for _, v := range vid {
			s := make(map[string]string)
			s["condition"] = "VocalId"
			s["keyword"] = v
			postjson := &vtbSearchJSON{
				Search:    s,
				PageIndex: 1,
				PageRows:  9999,
			}
			p, err := json.Marshal(postjson)
			if err != nil {
				log.Fatalln(err)
			}
			result, err := cqfunction.WebPostJSONContent(getAPIAddr("MusicList"), string(p))
			if err != nil {
				_, ok := err.(*cqfunction.TimeOutError)
				if ok {
					log.Println(err.Error())
					runtime.Goexit()
				} else {
					log.Fatalln(err)
				}
			}
			tmp := gjson.GetBytes(result, "Data").Array()
			for _, v2 := range tmp {
				ml.Data = append(ml.Data, v2)
			}
			ml.Total = ml.Total + len(ml.Data)
		}
	default:
		s := make(map[string]string)
		s["condition"] = "OriginName"
		s["keyword"] = Keyword
		postjson := &vtbSearchJSON{
			Search:    s,
			PageIndex: 1,
			PageRows:  9999,
		}

		p, err := json.Marshal(postjson)
		if err != nil {
			log.Fatalln(err)
		}

		result, err := cqfunction.WebPostJSONContent(getAPIAddr("MusicList"), string(p))
		if err != nil {
			_, ok := err.(*cqfunction.TimeOutError)
			if ok {
				log.Println(err.Error())
				runtime.Goexit()
			} else {
				log.Fatalln(err)
			}
		}
		ml.Total = int(gjson.GetBytes(result, "Total").Int())
		ml.Data = gjson.GetBytes(result, "Data").Array()
		return ml
	}
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
	result, err := cqfunction.WebPostJSONContent(getAPIAddr("CDN"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
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

//GetVTBsList : Get VTB Detail Info.
func GetVTBsList(VtbsName string) (VList *VtbsList) {
	vl := new(VtbsList)
	s := make(map[string]string)
	s["condition"] = "ChineseName"
	s["keyword"] = VtbsName
	postjson := vtbSearchJSON{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(getAPIAddr("VtbsList"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	vl.Data = gjson.GetBytes(result, "Data").Array()
	vl.Total = len(vl.Data)
	return vl
}

//GetVTBMusicDetail : Get music info by using music id.
func GetVTBMusicDetail(VTBMusicID string) (MusicInfo *MusicList) {
	ml := new(MusicList)
	postjson := vtbMusicData{
		MusicID: VTBMusicID,
	}
	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(getAPIAddr("MusicData"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
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

//GetHotMusicList : Get the hot music.
func GetHotMusicList() (VTBMusicList *MusicList) {
	ml := new(MusicList)
	postjson := vtbSearchJSON{
		PageIndex: 1,
		PageRows:  12,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(getAPIAddr("HotList"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	ml.Total = 12
	ml.Data = gjson.GetBytes(result, "Data").Array()
	return ml
}

type vtbCDNJSON struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
type vtbSearchJSON struct {
	Search    map[string]string `json:"search"`
	condition string
	keyword   string
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
type vtbMusicData struct {
	MusicID string `json:"id"`
}
type vtbHotMusic struct {
	PageIndex int `json:"pageIndex"`
	PageRows  int `json:"pageRows"`
}
