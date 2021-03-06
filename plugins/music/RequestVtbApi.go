package music

import (
	"encoding/json"
	"log"
	"runtime"

	"github.com/Ink-33/SMLKBOT/utils/cqfunction"
)

const vtbMusicAPI = "https://api.aqua.chat/v1/"

var vtbMusicAPIMap = make(map[string]string)

func init() {
	vtbMusicAPIMap["MusicList"] = "GetMusicList"
	vtbMusicAPIMap["MusicData"] = "GetMusicData"
	vtbMusicAPIMap["VtbsList"] = "GetVtbsList"
	vtbMusicAPIMap["VtbsData"] = "GetVtbsData"
	vtbMusicAPIMap["CDN"] = "GetCDNList"
	vtbMusicAPIMap["HotList"] = "GetHotMusicList"
}

func (e *VTBMusicClient) getAPIAddr(apiName string) string {
	return vtbMusicAPI + vtbMusicAPIMap[apiName]
}

// GetMusicList : Get VTBMusic Detail Info.
//	Method: MusicName,VtbName
func (e *VTBMusicClient) GetMusicList(keyword string, method string) (list *VTBMusicList, err error) {
	ml := new(VTBMusicList)
	s := make(map[string]string)
	s["keyword"] = keyword
	switch method {
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

	result, err := cqfunction.WebPostJSONContent(e.getAPIAddr("MusicList"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	decode := new(GetVTBMusicList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		log.Println(err)
		return ml, err
	}
	ml.Total = decode.Total
	ml.Data = decode.Data
	return ml, err
}

// GetMusicCDN : Get VTBMusic CDN Detail Info
func (e *VTBMusicClient) GetMusicCDN(keyword string) (cdn *GetVTBCDNList) {
	// log.Println(keyword)
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
	// log.Println(string(p))
	result, err := cqfunction.WebPostJSONContent(e.getAPIAddr("CDN"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	decode := new(GetVTBCDNList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		log.Println(err)
		return nil
	}
	return decode
}

// GetVTBsList : Get VTB Detail Info.
func (e *VTBMusicClient) GetVTBsList(vtbsName string) (vList *VTBsList) {
	vl := new(VTBsList)
	s := make(map[string]string)
	s["condition"] = "ChineseName"
	s["keyword"] = vtbsName
	postjson := vtbSearchJSON{
		Search:    s,
		PageIndex: 1,
		PageRows:  9999,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(e.getAPIAddr("VtbsList"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	decode := new(VTBsList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		log.Println(err)
		return nil
	}
	vl.Data = decode.Data
	vl.Total = len(vl.Data)
	return vl
}

// GetVTBMusicDetail : Get music info by using music id.
func (e *VTBMusicClient) GetVTBMusicDetail(vtbMusicID string) (musicInfo *VTBMusicList, err error) {
	ml := new(VTBMusicList)
	postjson := vtbMusicData{
		MusicID: vtbMusicID,
	}
	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
		ml.Total = -1
		return ml, err
	}

	result, err := cqfunction.WebPostJSONContent(e.getAPIAddr("MusicData"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	decode := new(GetVTBMusicData)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		return ml, err
	}
	if decode.GetVTBMusicListData != nil {
		dataArray := make([]GetVTBMusicListData, 0)
		dataArray = append(dataArray, *decode.GetVTBMusicListData)
		ml.Data = dataArray
		ml.Total = 1
		return ml, nil
	}
	ml.Total = 0
	return ml, nil
}

// GetHotMusicList : Get the hot music.
func (e *VTBMusicClient) GetHotMusicList() (list *VTBMusicList, err error) {
	ml := new(VTBMusicList)
	postjson := vtbSearchJSON{
		PageIndex: 1,
		PageRows:  12,
	}

	p, err := json.Marshal(postjson)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := cqfunction.WebPostJSONContent(e.getAPIAddr("HotList"), string(p))
	if err != nil {
		_, ok := err.(*cqfunction.TimeOutError)
		if ok {
			log.Println(err.Error())
			runtime.Goexit()
		} else {
			log.Fatalln(err)
		}
	}
	decode := new(GetVTBMusicList)
	err = json.Unmarshal(result, decode)
	if err != nil {
		ml.Total = -1
		return ml, err
	}
	ml.Total = decode.Total
	ml.Data = decode.Data
	return ml, nil
}

func (cdn *GetVTBCDNList) match(keyword string) (addr string) {
	if cdn.Data != nil {
		var addr string
		for i := range cdn.Data {
			if cdn.Data[i].Name == keyword {
				addr = cdn.Data[i].URL
				return addr
			}
		}
	}
	return ""
}

type vtbCDNJSON struct {
	Search    map[string]string `json:"search"`
	PageIndex int               `json:"pageIndex"`
	PageRows  int               `json:"pageRows"`
}

type vtbSearchJSON struct {
	Search    map[string]string `json:"search"`
	PageIndex int               `json:"pageIndex"`
	PageRows  int               `json:"pageRows"`
}

type vtbMusicData struct {
	MusicID string `json:"id"`
}
