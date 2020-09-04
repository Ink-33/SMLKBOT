package vtbmusic

import (
	"SMLKBOT/utils/cqfunction"
	"encoding/json"
	"log"
	"runtime"
)

/*
Reverse proxy for VtbmusicAPI: https://api.vtbmusic.com:60006/v1/
	This proxy is using Tencent Cloud API Gateway to record the usage.
*/
const vtbMusicAPIProxy string = "https://service-0pbekx7m-1252062863.gz.apigw.tencentcs.com/release/v1/"

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
	s := make(map[string]string)
	s["keyword"] = Keyword
	switch Method {
	case "VtbName":
		s["condition"] = "VocalName"
	default:
		s["condition"] = "OriginName"
	}
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
	decode := new(GetMusicList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		log.Println(err)
		return ml
	}
	ml.Total = decode.Total
	ml.Data = decode.Data
	return ml

}

//GetVTBMusicCDN : Get VTBMusic CDN Detail Info
func GetVTBMusicCDN(keyword string) (cdn *GetCDNList) {
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
	decode := new(GetCDNList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decode
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
	decode := new(GetVtbsList)
	err = json.Unmarshal(result, decode)
	vl.Data = decode.Data
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
	decode := new(GetMusicData)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		return ml
	}
	if decode.GetMusicListData != nil {
		dataArray := make([]GetMusicListData, 0)
		dataArray = append(dataArray, *decode.GetMusicListData)
		ml.Data = dataArray
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
	decode := new(GetMusicList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		return ml
	}
	ml.Total = decode.Total
	ml.Data = decode.Data
	return ml
}

func (cdn *GetCDNList) match(keyword string) (addr string) {
	if cdn.Data != nil {
		var addr string
		for _, r := range cdn.Data {
			if r.Name == keyword {
				addr = r.URL
				return addr
			}
		}
	}
	return ""
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
